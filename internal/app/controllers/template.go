package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.com/metronero/backend/internal/utils/config"
	"gitlab.com/metronero/backend/internal/utils/helpers"
	"gitlab.com/metronero/backend/pkg/apierror"
	"gitlab.com/metronero/backend/pkg/models"
)

// Get a preview of the payment page template
func GetMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountId := ctx.Value("account_id").(string)
	var t *template.Template
	p, err := url.JoinPath(config.TemplateDir, accountId, "index.html")
	if err != nil {
		helpers.WriteError(w, apierror.ErrTemplateLoad, err)
		return
	}
	t, err = template.ParseFiles(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t = config.DefaultPaymentTemplate
		} else {
			helpers.WriteError(w, apierror.ErrTemplateLoad, err)
			return
		}
	}

	// Execute with dummy data
	t.Execute(w, &models.PaymentPageInfo{
		InvoiceId:    "0b1b6a94-9ec2-4bdf-8251-6a46aca6a332",
		MerchantName: ctx.Value("username").(string),
		Amount:       120000000000,
		AmountFloat:  "0.12",
		OrderId:      "AI6X21",
		Status:       "Pending",
		LastUpdate:   time.Now(),
		Address:      "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw",
		ExtraData:    "Sneakers x 2, Jacket x 1",
		Qr:           "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACuElEQVR42uyZwY30KhCEy+LQR0IgE5PYaOyRE8OZEAJHDhb11M386519ARhLw2nl/Q4Dbrqq2viu77rhEpLcy8QXmR8AaqwzJn1abgRkQBI6sGSEqshUADcSQErC7F9cgaC0NNnLwm0wYGey010DMwmMCZSJDY5rOIAoe51uB+iumHzDIxwBQZJomX9W1NVAv3qzf5VH2PIRJNVZ9v/dzUsBWxP3fvWoR00y/Wl1FwOSEYUllqdfsVBr+t3CzndxPUBJNRarUh569VKdKnzzvBGgS2vDvxTJLteonbkAy0CACsTsd+sP5Ls/6D/cSIAkbcRlMuVV8SDr7Nu5iwEAkJh9wlIe0KtXI/c6sfntRoDQCiay+QPO3gVmScBQQFZXo3rhVZqPgMgm6XwXQwBkkgaYyXmAWY+aDRgKyNp7y2wFvgb1tFBPqyc/DtB3MXH1zC67TFYIy3LqxQ0ABM0Lfi9Pv3HFEVRA/rTiywGh2ZykGhaYD6jnYVNjhnGAjKg/2rM4C2JMNVZtF2cQuxzQJ2J/e3LLzEJNi+oXbgUAMGn2zND4ADbE8mk4LwYkv6UZ5ic1UFZIw/PsciMA+sTvdrLvVqwRuPkVAwG15xwWx43OANm5/lbe6wH2Kn3C5X7UwjrxVdydAOtxGs/Kkh1dRjQftv5SvesB7XIaadQyaDS3LqeRZx0HQBcIdWJHOIJjjX1Q86O8gwD/7GJe6HIfeZ31cAtAiChN/QOWbnohTauo/Gj39UCfgRSL5hnQXeirwdN/DGquBd7TJIUQNv5MF59nCBoAsGGyZ3niAUcVD2sYSzmd2ACATThtF3DZsYI7+Ttf3AVQLSaPsMLexVTxZ5tDABZx+zSJ1HSJqXx8XLga6F9AtPfaNEkoezWNGwj49y1JNWzhEUwvtMwfuA/wXd811PovAAD//y9+6+A69LGyAAAAAElFTkSuQmCC",
	})
}

func GetDefaultTemplatePreview(w http.ResponseWriter, r *http.Request) {
	config.DefaultPaymentTemplate.Execute(w, &models.PaymentPageInfo{
		InvoiceId:    "0b1b6a94-9ec2-4bdf-8251-6a46aca6a332",
		MerchantName: r.Context().Value("username").(string),
		Amount:       120000000000,
		AmountFloat:  "0.12",
		OrderId:      "AI6X21",
		Status:       "Pending",
		LastUpdate:   time.Now(),
		Address:      "46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw",
		ExtraData:    "Sneakers x 2, Jacket x 1",
		Qr:           "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAACuElEQVR42uyZwY30KhCEy+LQR0IgE5PYaOyRE8OZEAJHDhb11M386519ARhLw2nl/Q4Dbrqq2viu77rhEpLcy8QXmR8AaqwzJn1abgRkQBI6sGSEqshUADcSQErC7F9cgaC0NNnLwm0wYGey010DMwmMCZSJDY5rOIAoe51uB+iumHzDIxwBQZJomX9W1NVAv3qzf5VH2PIRJNVZ9v/dzUsBWxP3fvWoR00y/Wl1FwOSEYUllqdfsVBr+t3CzndxPUBJNRarUh569VKdKnzzvBGgS2vDvxTJLteonbkAy0CACsTsd+sP5Ls/6D/cSIAkbcRlMuVV8SDr7Nu5iwEAkJh9wlIe0KtXI/c6sfntRoDQCiay+QPO3gVmScBQQFZXo3rhVZqPgMgm6XwXQwBkkgaYyXmAWY+aDRgKyNp7y2wFvgb1tFBPqyc/DtB3MXH1zC67TFYIy3LqxQ0ABM0Lfi9Pv3HFEVRA/rTiywGh2ZykGhaYD6jnYVNjhnGAjKg/2rM4C2JMNVZtF2cQuxzQJ2J/e3LLzEJNi+oXbgUAMGn2zND4ADbE8mk4LwYkv6UZ5ic1UFZIw/PsciMA+sTvdrLvVqwRuPkVAwG15xwWx43OANm5/lbe6wH2Kn3C5X7UwjrxVdydAOtxGs/Kkh1dRjQftv5SvesB7XIaadQyaDS3LqeRZx0HQBcIdWJHOIJjjX1Q86O8gwD/7GJe6HIfeZ31cAtAiChN/QOWbnohTauo/Gj39UCfgRSL5hnQXeirwdN/DGquBd7TJIUQNv5MF59nCBoAsGGyZ3niAUcVD2sYSzmd2ACATThtF3DZsYI7+Ttf3AVQLSaPsMLexVTxZ5tDABZx+zSJ1HSJqXx8XLga6F9AtPfaNEkoezWNGwj49y1JNWzhEUwvtMwfuA/wXd811PovAAD//y9+6+A69LGyAAAAAElFTkSuQmCC",
	})
}

// Helper function to sanitize filenames
func sanitizeFileName(filename string) string {
	// Remove any path separators and unwanted characters
	filename = filepath.Base(filename)                 // Strips any path
	filename = strings.ReplaceAll(filename, "/", "_")  // Replace slashes
	filename = strings.ReplaceAll(filename, "\\", "_") // Replace backslashes

	// Optional: Use regex to allow only certain characters (alphanumeric, dashes, underscores, dots)
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	return reg.ReplaceAllString(filename, "_")
}

func logDirectoryContents(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		log.Info().
			Str("file", file.Name()).
			Str("directory", dirPath).
			Msg("File in directory")
	}
	return nil
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// Upload new templates.
func PostMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << config.TemplateMaxSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		helpers.WriteError(w, apierror.ErrBadRequest, err)
		return
	}

	ctx := r.Context()
	id, ok := ctx.Value("account_id").(string)
	if !ok {
		err := fmt.Errorf("account_id missing from context")
		log.Error().Err(err).Msg("Missing account ID")
		helpers.WriteError(w, apierror.ErrBadRequest, err)
		return
	}

	// Define the directory path for storing files for this account
	dirPath, err := url.JoinPath(config.TemplateDir, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create directory path")
		helpers.WriteError(w, apierror.ErrTemplateSave, err)
		return
	}

	// Ensure the directory exists or create it if it doesn't
	if err := ensureDir(dirPath); err != nil {
		log.Error().Err(err).Msg("Failed to create directory")
		helpers.WriteError(w, apierror.ErrTemplateSave, fmt.Errorf("failed to create directory: %w", err))
		return
	}

	// Get uploaded files under "template[]"
	files := r.MultipartForm.File["template[]"]
	if len(files) == 0 {
		err := fmt.Errorf("no files uploaded")
		log.Warn().Err(err).Msg("Upload attempt with no files")
		helpers.WriteError(w, apierror.ErrBadRequest, err)
		return
	}

	// Process each uploaded file
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			log.Error().Err(err).Str("filename", header.Filename).Msg("Failed to open file")
			helpers.WriteError(w, apierror.ErrBadRequest, err)
			return
		}
		defer file.Close()

		// Sanitize filename
		safeFilename := sanitizeFileName(header.Filename)

		// Create a unique path for each file in the directory
		filePath := filepath.Join(dirPath, safeFilename)

		// Create and open the destination file
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Error().Err(err).Str("file_path", filePath).Msg("Failed to open destination file")
			helpers.WriteError(w, apierror.ErrTemplateSave, err)
			return
		}
		defer destFile.Close()

		// Copy file data to the destination
		if _, err := io.Copy(destFile, file); err != nil {
			log.Error().Err(err).Str("file_path", filePath).Msg("Failed to save file")
			helpers.WriteError(w, apierror.ErrTemplateSave, err)
			return
		}

		log.Info().Str("file_path", filePath).Str("filename", safeFilename).Msg("File uploaded successfully")
	}

	// Log all files in the directory after uploads are complete
	if err := logDirectoryContents(dirPath); err != nil {
		log.Error().Err(err).Str("directory", dirPath).Msg("Failed to log directory contents")
	}

	// Respond with success
	log.Info().Str("directory", dirPath).Msg("All files uploaded successfully")
}

// Reset template back to default. Works by deleting merchant's template file.
func DeleteMerchantTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountId := ctx.Value("account_id").(string)
	p, err := url.JoinPath(config.TemplateDir, accountId)
	if err != nil {
		helpers.WriteError(w, apierror.ErrTemplateDelete, err)
		return
	}
	if err := os.RemoveAll(p); err != nil && !errors.Is(err, os.ErrNotExist) {
		helpers.WriteError(w, apierror.ErrTemplateDelete, err)
		return
	}
}
