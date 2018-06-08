// go-adorable - Adorable Avatars from Go
// Copyright (c) 2018 Péter Szilágyi. All rights reserved.

//go:generate go-bindata -nometadata -o assets.go -pkg bodyparts -ignore (((bodyparts)|(assets)).go|LICENSE) ./...
//go:generate gofmt -s -w assets.go

// Package bodyparts contains the avatar body parts.
//
// Courtesy of https://github.com/adorableio/avatars-api-middleware
package bodyparts
