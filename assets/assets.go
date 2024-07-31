package assets

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/izquiratops/bundler-tools/utils"

	"golang.org/x/net/html"
)

func handleAsset(srcDirPath, distDirPath, assetRelPath string) (string, error) {
	srcAbsPath, err := filepath.Abs(filepath.Join(srcDirPath, assetRelPath))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for source: %v", err)
	}

	distAbsPath, err := filepath.Abs(filepath.Join(distDirPath, filepath.Base(assetRelPath)))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for destination: %v", err)
	}

	if _, err := os.Stat(distDirPath); os.IsNotExist(err) {
		message := fmt.Sprintf("directory %s does not exist.\n", distDirPath)
		log.Fatal(message)
	}

	if err := utils.MoveTo(srcAbsPath, distAbsPath); err != nil {
		return "", err
	}

	// Setting the new src attribute as ./filename because everything is moved at the same path level
	return filepath.Base(assetRelPath), nil
}

func HandleHtmlFile(htmlBytes []byte, srcDirPath, distDirPath string) ([]byte, error) {
	document, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	var replaceAsset func(*html.Node) error
	replaceAsset = func(n *html.Node) error {
		// Handle SVG
		if n.Type == html.ElementNode && n.Data == "img" {
			for i, attr := range n.Attr {
				if attr.Key == "src" && bytes.HasSuffix([]byte(attr.Val), []byte(".svg")) {
					n.Attr[i].Val, err = handleAsset(srcDirPath, distDirPath, attr.Val)
					if err != nil {
						return err
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := replaceAsset(c); err != nil {
				return err
			}
		}

		return nil
	}

	if err := replaceAsset(document); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = html.Render(&buf, document); err != nil {
		return nil, fmt.Errorf("error rendering HTML: %w", err)
	}

	return buf.Bytes(), nil
}
