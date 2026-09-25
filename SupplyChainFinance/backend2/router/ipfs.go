package router

import (
	"backend/dbs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

func FileUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析表单数据: " + err.Error()})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有文件上传"})
		return
	}

	tag := strings.TrimSpace(c.PostForm("tag"))
	operator := strings.TrimSpace(c.PostForm("operator"))
	receivableID := strings.TrimSpace(c.PostForm("receivableId"))
	orderID := strings.TrimSpace(c.PostForm("orderId"))
	docType := strings.TrimSpace(c.PostForm("docType"))
	status := strings.TrimSpace(c.PostForm("status"))
	source := strings.TrimSpace(c.PostForm("source"))
	if docType == "" {
		docType = "other"
	}
	if status == "" {
		status = "uploaded"
	}
	if source == "" {
		source = "ipfs"
	}

	saveDir := "./uploadsIPFS"
	if err = os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建上传目录: " + err.Error()})
		return
	}

	cids := make([]string, 0, len(files))
	for _, file := range files {
		filePath := filepath.Join(saveDir, file.Filename)
		src, openErr := file.Open()
		if openErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件: " + openErr.Error()})
			return
		}

		dst, createErr := os.Create(filePath)
		if createErr != nil {
			_ = src.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建临时文件: " + createErr.Error()})
			return
		}

		_, copyErr := io.Copy(dst, src)
		_ = src.Close()
		_ = dst.Close()
		if copyErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存文件: " + copyErr.Error()})
			return
		}

		cid, uploadErr := uploadToIPFS(filePath)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法上传到 IPFS: " + uploadErr.Error()})
			return
		}

		if insertErr := dbs.InsertDocument(dbs.DocumentDB{
			Name:         file.Filename,
			CID:          cid,
			ReceivableID: receivableID,
			OrderID:      orderID,
			DocType:      docType,
			Status:       status,
			Source:       source,
			Tag:          tag,
			Operator:     operator,
		}); insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "IPFS 已上传，但写入数据库失败: " + insertErr.Error(), "cid": cid})
			return
		}

		cids = append(cids, cid)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件已上传至 IPFS，并同步写入数据库",
		"cids":    cids,
	})
}

func uploadToIPFS(filePath string) (string, error) {
	sh := shell.NewShell("localhost:5001")
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return sh.Add(f)
}
