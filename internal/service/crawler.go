package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/paulmach/orb"
	"github.com/xuri/excelize/v2"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

var (
	coordRe = regexp.MustCompile(`(\d{4,6})([NS])(\d{5,7})([EW])`) // ddmm(ss)Ndddmm(ss)E
	timeRe  = regexp.MustCompile(`^\d{4}$`)                        // ччмм
	dateRe  = regexp.MustCompile(`^\d{6}$`)                        // ггммдд
)

type ParserService struct {
	repo Repository
}

func NewParserService(repo Repository) *ParserService {
	return &ParserService{repo: repo}
}

func (p *ParserService) cleanString(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

func (p *ParserService) ProcessXLSX(ctx context.Context, f *excelize.File, authorID, filename string) (int, int, error) {

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Printf("Error reading rows from %s: %v", filename, err)
		return 0, 0, err
	}

	seen := make(map[string]struct{})
	reportLines := []string{}
	validCount := 0
	errorCount := 0
	for i, row := range rows {
		if len(row) < 2 {
			errorCount++
			continue
		}

		region := p.cleanString(row[0])
		shrRaw := p.cleanString(row[1])
		idepRaw := ""
		iarrRaw := ""
		if len(row) > 2 {
			idepRaw = p.cleanString(row[2])
		}
		if len(row) > 3 {
			iarrRaw = p.cleanString(row[3])
		}

		if region == "" || shrRaw == "" {
			errorCount++
			continue
		}

		if i < 5 {
			reportLines = append(reportLines, fmt.Sprintf("Row %d debug - Region: '%s', SHR: '%s'", i+1, region, shrRaw))
		}

		msg, _, _ := p.parseSHR(shrRaw, region)

		if idepRaw != "" {
		}
		if iarrRaw != "" {
			ata := p.parseATAFromIARR(iarrRaw)
			if ata != "" {
				msg.ATA = ata
			}
		}

		key := msg.SID + msg.DOF + msg.ATD
		if _, exists := seen[key]; exists {
			reportLines = append(reportLines, fmt.Sprintf("Row %d duplicate: %s", i+1, key))
			errorCount++
			continue
		}
		seen[key] = struct{}{}
		if msg.SID == "" || msg.DepCoords == "" || msg.ATA == "" {
			errorCount++
			continue
		}
		err = p.repo.SaveMessage(context.Background(), &msg)
		if err != nil {
			slog.Error("error saving message:", err.Error())
			errorCount++

			continue
		}
		validCount++
	}

	return validCount, errorCount, nil
}

func (p *ParserService) parseSHR(raw string, region string) (model.ParsedMessage, []string, []string) {
	msg := model.ParsedMessage{Region: region}
	changes := []string{}
	errs := []string{}

	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, `()"`)
	raw = strings.ReplaceAll(raw, " ", "")

	changes = append(changes, fmt.Sprintf("Raw SHR input: '%s'", raw))

	parts := strings.Split(raw, "-")
	changes = append(changes, fmt.Sprintf("SHR parts count: %d", len(parts)))

	if len(parts) < 3 {
		errs = append(errs, "Invalid SHR format - too few parts")
		return msg, changes, errs
	}

	if !strings.Contains(strings.ToUpper(parts[0]), "SHR") {
		errs = append(errs, "Invalid SHR format - missing SHR prefix")
		return msg, changes, errs
	}

	msg.REG = parts[1]
	changes = append(changes, fmt.Sprintf("Set REG from index: %s", msg.REG))

	if len(parts) > 2 {
		atdPart := parts[2]
		changes = append(changes, fmt.Sprintf("ATD part: '%s'", atdPart))

		if strings.HasPrefix(atdPart, "ZZZZ") {
			atd := atdPart[4:]
			changes = append(changes, fmt.Sprintf("ATD time: '%s'", atd))

			if valid, ch := p.validateTime(atd); valid {
				msg.ATD = ch
				changes = append(changes, fmt.Sprintf("Valid ATD: %s", ch))
			} else {
				errs = append(errs, fmt.Sprintf("Invalid ATD: '%s'", atd))
			}
		} else {
			if len(atdPart) >= 4 {
				if valid, ch := p.validateTime(atdPart); valid {
					msg.ATD = ch
					changes = append(changes, fmt.Sprintf("Found ATD without ZZZZ prefix: %s", ch))
				} else {
					errs = append(errs, fmt.Sprintf("ATD part doesn't start with ZZZZ and not valid time: '%s'", atdPart))
				}
			} else {
				errs = append(errs, fmt.Sprintf("ATD part too short: '%s'", atdPart))
			}
		}
	} else {
		errs = append(errs, "No ATD part found")
	}

	minAlt := 0
	maxAlt := 0
	foundHeight := false
	heightRe := regexp.MustCompile(`(K\d{4})?M(\d{4})(/M(\d{4}))?`)
	for i := 3; i < len(parts) && !foundHeight; i++ {
		matches := heightRe.FindStringSubmatch(parts[i])
		if len(matches) > 0 {
			foundHeight = true
			maxStr := matches[2]
			max, err := strconv.Atoi(maxStr)
			if err == nil {
				maxAlt = max * 10
			} else {
				errs = append(errs, fmt.Sprintf("Invalid max height: %s", maxStr))
			}

			if matches[3] != "" && matches[4] != "" {
				minStr := matches[2]
				min, err := strconv.Atoi(minStr)
				if err == nil {
					minAlt = min * 10
				} else {
					errs = append(errs, fmt.Sprintf("Invalid min height: %s", minStr))
				}
				maxStr = matches[4]
				max, err = strconv.Atoi(maxStr)
				if err == nil {
					maxAlt = max * 10
				} else {
					errs = append(errs, fmt.Sprintf("Invalid max height: %s", maxStr))
				}
			} else {
				minAlt = 0
			}
			changes = append(changes, fmt.Sprintf("Parsed heights: min=%d m, max=%d m", minAlt, maxAlt))
		}
	}
	if !foundHeight {
		changes = append(changes, "No height information found")
	}
	msg.MinAlt = minAlt
	msg.MaxAlt = maxAlt

	// Поле 15: высоты /ZONA ... - парсим зону полета
	zoneCoords, zoneLatLon := p.parseZone(parts, changes)
	msg.ZoneCoords = zoneCoords
	msg.ZoneLatLon = zoneLatLon

	// Поле 16: ZZZZччмм (EET)
	if len(parts) > 4 {
		eetPart := parts[4]
		if strings.HasPrefix(eetPart, "ZZZZ") {
			eet := eetPart[4:]
			if valid, ch := p.validateTime(eet); valid {
				msg.ATA = ch
				changes = append(changes, fmt.Sprintf("Added ATA from EET: %s", ch))
			}
		}
	}

	var field18 string
	if len(parts) > 5 {
		field18 = strings.Join(parts[5:], "-")
	} else if len(parts) > 4 {
		field18 = parts[4]
	}

	changes = append(changes, fmt.Sprintf("Field 18: '%s'", field18))
	kv := p.parseField18(field18)
	changes = append(changes, fmt.Sprintf("Parsed fields: %v", kv))

	if dof, ok := kv["DOF/"]; ok {
		if valid, ch := p.validateDate(dof); valid {
			msg.DOF = ch
			changes = append(changes, fmt.Sprintf("Normalized DOF: %s -> %s", dof, ch))
		} else {
			errs = append(errs, "Invalid DOF")
		}
	}
	if dep, ok := kv["DEP/"]; ok {
		changes = append(changes, fmt.Sprintf("Found DEP: '%s'", dep))
		if valid, norm, latlon := p.validateCoords(dep); valid {
			msg.DepCoords = norm
			msg.DepLatLon = latlon
			changes = append(changes, fmt.Sprintf("Normalized DEP: %s -> %s", dep, norm))
		} else {
			errs = append(errs, fmt.Sprintf("Invalid DEP coords: '%s'", dep))
		}
	} else {
		changes = append(changes, "No DEP field found")
	}

	if dest, ok := kv["DEST/"]; ok {
		changes = append(changes, fmt.Sprintf("Found DEST: '%s'", dest))
		if valid, norm, latlon := p.validateCoords(dest); valid {
			msg.ArrCoords = norm
			msg.ArrLatLon = latlon
			changes = append(changes, fmt.Sprintf("Normalized DEST: %s -> %s", dest, norm))
		} else {
			errs = append(errs, fmt.Sprintf("Invalid DEST coords: '%s'", dest))
		}
	} else {
		changes = append(changes, "No DEST field found")
	}
	if sid, ok := kv["SID/"]; ok {
		msg.SID = sid
		changes = append(changes, fmt.Sprintf("Found SID: %s", sid))
	} else {
		changes = append(changes, "No SID field found")
	}
	msg.OPR = kv["OPR/"]
	if reg, ok := kv["REG/"]; ok {
		msg.REG = reg
	}
	msg.TYP = kv["TYP/"]
	msg.RMK = kv["RMK/"]

	for k := range kv {
		if !strings.HasPrefix(k, "DEP/") && !strings.HasPrefix(k, "DEST/") && !strings.HasPrefix(k, "DOF/") && !strings.HasPrefix(k, "SID/") && !strings.HasPrefix(k, "OPR/") && !strings.HasPrefix(k, "REG/") && !strings.HasPrefix(k, "TYP/") && !strings.HasPrefix(k, "RMK/") && !strings.HasPrefix(k, "STS/") && !strings.HasPrefix(k, "EET/") {
			changes = append(changes, fmt.Sprintf("Ignored extra field: %s", k))
		}
	}

	return msg, changes, errs
}

func (p *ParserService) parseField18(raw string) map[string]string {
	kv := make(map[string]string)

	if raw == "" {
		return kv
	}

	patterns := map[string]string{
		"DOF/":  `DOF/(\d{6})`,
		"DEP/":  `DEP/([0-9]+[NS][0-9]+[EW])`,
		"DEST/": `DEST/([0-9]+[NS][0-9]+[EW])`,
		"SID/":  `SID/(\d+)`,
		"OPR/":  `OPR/([^A-Z/]+)`,
		"REG/":  `REG/([^A-Z/]+)`,
		"TYP/":  `TYP/([A-Z]+)`,
		"RMK/":  `RMK/([^A-Z/]+)`,
	}

	for key, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(raw); len(match) > 1 {
			kv[key] = strings.TrimSpace(match[1])
		}
	}

	return kv
}

// парсинг зоны полета из поля 15

func (p *ParserService) parseZone(parts []string, changes []string) ([]string, []orb.Point) {
	var zoneCoords []string
	var zoneLatLon []orb.Point

	// ищем поле с /ZONA
	for _, part := range parts {
		if strings.Contains(part, "/ZONA") {
			changes = append(changes, fmt.Sprintf("Found ZONA field: %s", part))

			// извлекаем координаты из зоны
			zonePattern := `([0-9]+[NS][0-9]+[EW])`
			re := regexp.MustCompile(zonePattern)
			matches := re.FindAllString(part, -1)

			for _, coord := range matches {
				zoneCoords = append(zoneCoords, coord)

				// конвертируем в decimal
				if valid, norm, latlon := p.validateCoords(coord); valid {
					zoneLatLon = append(zoneLatLon, latlon)
					changes = append(changes, fmt.Sprintf("Zone coord: %s -> %s", coord, norm))
				} else {
					changes = append(changes, fmt.Sprintf("Invalid zone coord: %s", coord))
				}
			}

			changes = append(changes, fmt.Sprintf("Found %d zone coordinates", len(zoneCoords)))
			break
		}
	}

	return zoneCoords, zoneLatLon
}

// валидация и нормализация времени ччмм -> чч:мм

func (p *ParserService) validateTime(t string) (bool, string) {
	t = strings.TrimSpace(t)
	if !timeRe.MatchString(t) {
		return false, ""
	}
	if len(t) != 4 {
		return false, ""
	}

	h, err1 := strconv.Atoi(t[:2])
	m, err2 := strconv.Atoi(t[2:])

	if err1 != nil || err2 != nil {
		return false, ""
	}

	if h < 0 || h > 23 || m < 0 || m > 59 {
		return false, ""
	}

	return true, fmt.Sprintf("%02d:%02d", h, m)
}

// валидация и нормализация даты

func (p *ParserService) validateDate(d string) (bool, string) {
	if !dateRe.MatchString(d) {
		return false, ""
	}
	y := "20" + d[:2]
	m := d[2:4]
	dd := d[4:]
	_, err := time.Parse("2006-01-02", y+"-"+m+"-"+dd)
	if err != nil {
		return false, ""
	}
	return true, y + "-" + m + "-" + dd
}

// валидация и нормализация координат

func (p *ParserService) validateCoords(c string) (bool, string, orb.Point) {
	m := coordRe.FindStringSubmatch(c)
	if len(m) != 5 {
		return false, "", orb.Point{}
	}
	latStr, ns, lonStr, ew := m[1], m[2], m[3], m[4]

	// Нормализация к ddmmss / dddmmss
	if len(latStr) == 4 {
		latStr += "00"
	} else if len(latStr) == 5 {
		latStr = latStr[:4] + "0" + latStr[4:]
	}
	if len(lonStr) == 5 {
		lonStr += "00"
	} else if len(lonStr) == 6 {
		lonStr = lonStr[:5] + "0" + lonStr[5:]
	}

	// decimal
	latD := float64(p.strToInt(latStr[:2]))
	latM := float64(p.strToInt(latStr[2:4]))
	latS := float64(p.strToInt(latStr[4:6]))
	lat := latD + latM/60 + latS/3600
	if ns == "S" {
		lat = -lat
	}
	if lat < -90 || lat > 90 {
		return false, "", [2]float64{}
	}

	lonD := float64(p.strToInt(lonStr[:3]))
	lonM := float64(p.strToInt(lonStr[3:5]))
	lonS := float64(p.strToInt(lonStr[5:7]))
	lon := lonD + lonM/60 + lonS/3600
	if ew == "W" {
		lon = -lon
	}
	if math.Abs(lon) > 180 {
		return false, "", [2]float64{}
	}

	norm := latStr + ns + lonStr + ew
	return true, norm, orb.Point{lon, lat}
}

func (p *ParserService) strToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (p *ParserService) parseATAFromIARR(raw string) string {
	parts := strings.Split(raw, "-")
	for _, part := range parts {
		if strings.HasPrefix(part, "ATA ") {
			t := strings.TrimSpace(part[4:])
			if valid, norm := p.validateTime(t); valid {
				return norm
			}
		}
	}
	return ""
}
