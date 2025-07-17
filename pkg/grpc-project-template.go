package pkg

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github/weeback/grpc-project-template/pkg/net"
)

const (
	OriginalURL = "https://github/weeback/grpc-project-template.git"
	LicenseURL  = "https://github/weeback/grpc-project-template/blob/main/LICENSE"

	Author  = "Bankaool, S.A., Institución de Banca Múltiple"
	License = "MIT"

	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
)

var (
	Version     = "0.0.1"
	BuildDate   = "2025-04-12T00:00:00Z"
	BuildUser   = "Unknown"
	BuildBranch = "Unknown"
	BuildCommit = "[?]"
	BuildTag    string // default is empty

	RepoURL     string
	LicenseText = `Restricted Use License

Copyright (c) 2025 Bankaool, S.A., Institución de Banca Múltiple

Permission is hereby granted to the designated organization to obtain a copy
of this software and associated documentation files (the "Software"), to use
within the organization without restriction, including but not limited to the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software within the organization, and to permit persons within the organization to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
)

// On import package 'github/weeback/grpc-project-template' to use, this function will be called.
func init() {
	if BuildTag == "" {
		Version = fmt.Sprintf("%s | Branch: %s - Commit: %s", Version, BuildBranch, BuildCommit)
	} else {
		Version = BuildTag
	}
	if RepoURL == "" {
		RepoURL = OriginalURL
	}

	println(ColorYellow)
	println(GetServiceInfo())
	println("-----------------------------------------------------------------------")
	println(ColorGreen)
	println(GetBuiltVersionInfo())
	println("=======================================================================")
	println(ColorReset)
}

func GetServiceInfo() string {
	return fmt.Sprintf("=== Business Center Service ===\n"+
		"Author: %s\nLicense: %s\nRepository: %s", Author, LicenseURL, RepoURL)
}

func GetBuiltVersionInfo() string {
	return fmt.Sprintf(">> Version: %s\n>> Build by %s\n>> Build at %s",
		Version, BuildUser, BuildDate)
}

// Import to validate LICENSE file exist, if not create.
func Import() {
	if _, err := os.Stat("LICENSE"); os.IsNotExist(err) {
		if err := os.WriteFile("LICENSE", []byte(LicenseText), 0644); err != nil {
			println("Error: Failed to create LICENSE file.")
		}
	}
}

func PermissionDeniedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := struct{}{}
		payload := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Permission Denied</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 40px;
			background-color: #f4f4f4;
			color: #333;
		}
		h1 {
			color: #d9534f;
			border-bottom: 2px solid #d9534f;
			padding-bottom: 5px;
		}
		p {
	
			font-size: 18px;
			line-height: 1.6;
		}
	
		a {
			color: #337ab7;
			text-decoration: none;
		}
		a:hover {
			text-decoration: underline;
		}
	
	
	</style>
</head>
<body>
	
	<h1>Permission Denied</h1>
	<p>You do not have permission to access this resource.</p>
	
	<p>If you believe this is an error, please contact the system administrator.</p>
	
	
	<p>For more information, please refer to the <a href="https://github/weeback/grpc-project-template/blob/main/LICENSE" target="_blank">License</a>.</p>
	
</body>
</html>`

		// parse template with param variables
		temp, err := template.New("healthcheck").Parse(payload)
		if err != nil {
			net.WriteError(w, http.StatusServiceUnavailable, err)
			return
		}
		// write status header
		w.WriteHeader(http.StatusOK)
		// write html file
		if err := temp.Execute(w, &params); err != nil {
			net.WriteError(w, http.StatusServiceUnavailable, err)
			return
		}
	})
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Author                 string
		ProtoURL, ProtoUrlText string
		Version                string
		Builder, BuildDate     string
	}{
		Author:       Author,
		ProtoURL:     "/proto/",
		ProtoUrlText: "here",
		Version:      Version,
		Builder:      BuildUser,
		BuildDate:    BuildDate,
	}
	payload := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Service Information</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 40px;
        }
        h1 {
            border-bottom: 2px solid #333;
            padding-bottom: 5px;
        }
        .section {
            margin-top: 30px;
        }
        .info-table {
            border-collapse: collapse;
            width: 100%;
            max-width: 600px;
        }
        .info-table th, .info-table td {
            text-align: left;
            padding: 8px;
        }
        .info-table th {
            background-color: #f2f2f2;
            width: 150px;
        }
        .info-table td {
            background-color: #fafafa;
        }
        a {
            color: #007acc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>

    <h1>Service Information</h1>

    <div class="section">
        <h2>GetService</h2>
        <table class="info-table">
            <tr>
                <th>Author</th>
                <td>{{.Author}}</td> <!-- Thay thế bằng dữ liệu thực tế -->
            </tr>
            <tr>
                <th>gRPC Proto URL</th>
                <td><a href="{{.ProtoURL}}" target="_blank">{{.ProtoUrlText}}</a></td>
            </tr>
        </table>
    </div>

    <div class="section">
        <h2>BuiltVersion</h2>
        <table class="info-table">
            <tr>
                <th>Version</th>
                <td>{{.Version}}</td> <!-- Thay thế bằng version thực tế -->
            </tr>
            <tr>
                <th>Build by</th>
                <td>{{.Builder}}</td> <!-- Tên người hoặc hệ thống build -->
            </tr>
            <tr>
                <th>Build at</th>
                <td>{{.BuildDate}}</td> <!-- Thời gian build -->
            </tr>
        </table>
    </div>

</body>
</html>`

	// parse template with param variables
	temp, err := template.New("healthcheck").Parse(payload)
	if err != nil {
		net.WriteError(w, http.StatusServiceUnavailable, err)
		return
	}
	// write status header
	w.WriteHeader(http.StatusOK)
	// write html file
	if err := temp.Execute(w, &params); err != nil {
		net.WriteError(w, http.StatusServiceUnavailable, err)
		return
	}
}

// GetFullPathFromRoot to call os.Getwd from root working directory.
func GetFullPathFromRoot(path string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, path)
}
