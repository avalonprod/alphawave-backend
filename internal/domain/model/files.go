package model

import "time"

const (
	BUCKET_STORAGE = "storage"
	BUCKET_IMAGES  = "images"
)

type FolderPath struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type Folder struct {
	ID               string       `json:"id" bson:"_id,omitempty"`
	TeamId           string       `json:"teamId" bson:"teamId"`
	Name             string       `json:"name" bson:"name"`
	Type             string       `json:"type" bson:"type"`
	Path             []FolderPath `json:"path" bson:"path"`
	CreatedAt        time.Time    `json:"createdAt" bson:"createdAt"`
	LastModifiedTime time.Time    `json:"lastModifiedTime" bson:"lastModifiedTime"`
	ParentFolder     string       `json:"parentFolder" bson:"parentFolder"`
}

type File struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	TeamId    string    `json:"teamId" bson:"teamId"`
	OwnerName string    `json:"ownerName" bson:"ownerName"`
	Name      string    `json:"name" bson:"name"`
	FilePath  string    `json:"filePath" bson:"filePath"`
	FolderId  string    `json:"folderId" bson:"folderId"`
	Url       string    `json:"url" bson:"url"`
	Key       string    `json:"key" bson:"key"`
	Type      string    `json:"type" bson:"type"`
	Path      []string  `json:"path" bson:"path"`
	Size      int       `json:"size" bson:"size"`
	Extension string    `json:"extension" bson:"extension"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
