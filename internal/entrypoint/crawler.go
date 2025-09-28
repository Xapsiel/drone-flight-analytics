package httpv1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"

	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

// UploadFileHandler
// @Summary Загрузить XLSX с данными
// @Description Принимает файл .xlsx, запускает парсинг и обновление метрик
// @Tags crawler
// @Accept mpfd
// @Produce json
// @Param file formData file true "Файл XLSX"
// @Param authorID formData string true "Идентификатор автора"
// @Success 200 {object} httpv1.APIResponse
// @Failure 400 {object} httpv1.APIResponse
// @Failure 500 {object} httpv1.APIResponse
// @Router /crawler/upload [post]
func (r *Router) UploadFileHandler(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		slog.Error("failed to read uploaded file", "error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Ошибка загрузки файла"))
	}
	filename := file.Filename
	if filename == "" || !strings.HasSuffix(filename, ".xlsx") {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "invalid file extension"))
	}
	authorID := ctx.FormValue("authorID")
	if authorID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Отсутствует authorID"))
	}
	fileChan := make(chan int)
	go func(file *multipart.FileHeader, ch chan int) {

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
		file_id, err := r.repo.SaveFileInfo(context.Background(), mf, 0, 0)
		ch <- file_id
		if err != nil {
			slog.Error(fmt.Sprintf("error saving file: %v", err))
		}

		validCount, errorCount, parseErr := r.service.ParserService.ProcessXLSX(context.Background(), f, authorID, file.Filename, file_id)
		if parseErr != nil {
			slog.Error("failed to parse xlsx", "filename", filename, "error", parseErr)
			return
		}
		mf.Status = "parsed"
		_, err = r.repo.SaveFileInfo(context.Background(), mf, validCount, errorCount)
		if err != nil {
			slog.Error(fmt.Sprintf("error saving file: %v", err))
			return
		}
		err = r.service.MetricsService.Update(context.Background())
		if err != nil {
			slog.Error(fmt.Sprintf("error update metrics: %v", err))
			return
		}
		slog.Info(fmt.Sprintf("successfully processed file: %v", mf.Filename))

	}(file, fileChan)
	fileID := <-fileChan
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(fiber.Map{
		"message":  "Файл успешно обработан",
		"authorID": authorID,
		"filename": file.Filename,
		"size":     file.Size,
		"file_id":  fileID,
	}, ""))
}

func (r *Router) CheckFileStatus(ctx *fiber.Ctx) error {
	fileID, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		slog.Error("failed to parse file id", "id", ctx.Query("id"))
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Отсутствует fileID"))
	}
	f, err := r.repo.GetFile(context.Background(), fileID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(r.NewErrorResponse(fiber.StatusBadRequest, "Ошибка с поиском файла"))
	}
	return ctx.Status(fiber.StatusOK).JSON(r.NewSuccessResponse(
		f, ""))

}
