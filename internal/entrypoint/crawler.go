package httpv1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

func (r *Router) UploadFileHandler(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	filename := file.Filename
	if filename == "" || !strings.HasSuffix(filename, ".xlsx") {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
		})
	}
	authorID := ctx.FormValue("authorID")
	if authorID == "" {
		return fmt.Errorf("authorID is required")
	}

	go func(file *multipart.FileHeader) {
		jsonData, err := json.Marshal(file.Header)
		if err != nil {
			slog.Error(fmt.Sprintf("error with marshaling: %v", err))
		}
		mf := model.File{
			AuthorID: authorID,
			Filename: filename,
			Size:     file.Size,
			Metadata: jsonData,
			Status:   "processing",
		}
		fi, err := file.Open()
		if err != nil {
			slog.Error(fmt.Sprintf("error opening file: %v", err))
			return
		}
		defer fi.Close()

		f, err := excelize.OpenReader(fi)
		if err != nil {
			slog.Error(fmt.Sprintf("error opening file: %v", err))
			return
		}
		err = r.repo.SaveFileInfo(context.Background(), mf, 0, 0)
		if err != nil {
			slog.Error(fmt.Sprintf("error saving file: %v", err))
		}

		validCount, errorCount, err := r.service.ParserService.ProcessXLSX(context.Background(), f, authorID, file.Filename)
		mf.Status = "parsed"
		err = r.repo.SaveFileInfo(context.Background(), mf, validCount, errorCount)
		if err != nil {
			slog.Error(fmt.Sprintf("error saving file: %v", err))
			return
		}
	}(file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка обработки XLSX: " + err.Error()})
	}
	return ctx.JSON(fiber.Map{
		"message":  "Файл успешно обработан",
		"authorID": authorID,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
