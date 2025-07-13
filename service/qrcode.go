package services

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// QRCodeService 二维码服务
type QRCodeService struct {
	qrCodeDir string
}

// NewQRCodeService 创建二维码服务实例
func NewQRCodeService(qrCodeDir string) *QRCodeService {
	return &QRCodeService{
		qrCodeDir: qrCodeDir,
	}
}

// GenerateQRCode 生成二维码图片
func (s *QRCodeService) GenerateQRCode(content string, size int) (string, error) {
	// 创建二维码目录（如果不存在）
	if err := os.MkdirAll(s.qrCodeDir, 0755); err != nil {
		return "", err
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("qr_%d.png", time.Now().UnixNano())
	filePath := filepath.Join(s.qrCodeDir, fileName)

	// 生成二维码
	qrCode, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return "", err
	}

	// 调整大小
	qrCode, err = barcode.Scale(qrCode, size, size)
	if err != nil {
		return "", err
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 保存二维码图片
	err = png.Encode(file, qrCode)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// GetQRCodePath 获取二维码文件路径
func (s *QRCodeService) GetQRCodePath(fileName string) string {
	return filepath.Join(s.qrCodeDir, fileName)
}
