package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"time"
)

type WebSocketURL struct {
	WebSocket string `json:"webSocketDebuggerUrl"`
}

func InitWebSocket() (string, error) {
	KillExistingChrome()

	chromePath := `C:\Program Files\Google\Chrome\Application\chrome.exe`
    userDataDir := `C:\ChromeDebug`

    cmd := exec.Command(
        chromePath,
        "--remote-debugging-port=9223",
        "--user-data-dir="+userDataDir,
        "--disable-gpu",
        "--no-first-run",
    )

	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("Failed open chrome.exe for port=9223: %v", err)
	}

	if !waitForPort(9223, 30*time.Second) {
        cmd.Process.Kill() // Завершаем процесс если порт не открылся
        return "", fmt.Errorf("chrome failed to start within timeout")
    }

	var websocket WebSocketURL

	err = getJSON(&websocket, "http://localhost:9223/json")
	if err != nil {
		return "", fmt.Errorf("failed get websocket: %v", err)
	}

	return websocket.WebSocket, nil
}

func waitForPort(port int, timeout time.Duration) bool {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 1*time.Second)
        if err == nil {
            conn.Close()
            time.Sleep(2 * time.Second) // Дополнительная задержка для полного запуска
            return true
        }
        time.Sleep(500 * time.Millisecond)
    }
    return false
}

func getJSON(websocket *WebSocketURL, path string) (error) {
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var targets []WebSocketURL
    if err := json.Unmarshal(body, &targets); err != nil {
        return err
    }

    if len(targets) == 0 {
        return fmt.Errorf("no debug targets found")
    }

    *websocket = targets[0]
	
	return nil
}

func KillExistingChrome() {
    // Можно добавить логику для завершения предыдущих процессов Chrome
    exec.Command("taskkill", "/F", "/IM", "chrome.exe").Run()
}