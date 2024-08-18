package entities

type FileMetadata struct {
	FileName        string `json:"file_name"`
	FileExtension   string `json:"file_extension"`
	MemType         string `json:"mem_type"`
	IsCompressed    bool   `json:"is_compressed"`
	CompressionType string `json:"compression_type"`
	FileSize        int64  `json:"file_size"`
}
