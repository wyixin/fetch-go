package dom

import (
	"fmt"
	"testing"
)

func getTestContent() string {
	content := `
<!DOCTYPE html>
<html>
<head>
    <title>Sample HTML Page</title>
    <link rel="stylesheet" href="http://baidu.com/styles.css">
</head>
<body>
    <h1>Welcome to My Sample HTML Page</h1>
    <p>This is a simple HTML page with an image and external CSS.</p>
    <img src="http://baidu.com/static/sample-image.jpg" alt="Sample Image">
</body>
</html>
`
	return content
}

func TestParseAllAssets(t *testing.T) {
	content := getTestContent()
	jsfiles, imgfiles, cssfiles := ParseAllAssets(content)
	fmt.Println(jsfiles, imgfiles, cssfiles)
}
