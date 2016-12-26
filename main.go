package main
 
import (
    "log"
    "net"
    "fmt"
    "os"
    "os/exec"
    "bytes"
    "runtime"
)
 
func main() {
    SHUTDOWN_PROGRAM := "shutdown"
    MAGIC_PACKET := []byte("shutdownd|trigger|CHANGEME\n")
    SUCCESS_MSG := []byte("shutdownd|ok\n")

    for _, v := range os.Args {
        switch v {
            case "-D":
                SHUTDOWN_PROGRAM = "echo"
                log.Printf("Dummy mode")
            default:
                MAGIC_PACKET = []byte(fmt.Sprintf("shutdownd|trigger|%s\n", v))
        }
    }

    sAddr, err := net.ResolveUDPAddr("udp", ":10001")
    if err != nil {
        panic(err)
    }
 
    sConn, err := net.ListenUDP("udp", sAddr)
    if err != nil {
        panic(err)
    }
    defer sConn.Close()
 
    buf := make([]byte, 1024)
 
    for {
        n, addr, err := sConn.ReadFromUDP(buf)
        if err != nil {
            log.Printf("Error reading: %v", err)
            continue
        }

        if bytes.Compare(buf[:n], MAGIC_PACKET) != 0 {
            log.Printf("Invalid packet from: %v", addr)
            continue
        }

        log.Printf("Shutdown trigger from: %v", addr)

        var cmd *exec.Cmd
        switch os := runtime.GOOS; os {
        case "windows":
            cmd = exec.Command(SHUTDOWN_PROGRAM, "-s", "-f", "-t", "60")
        case "darwin":
            cmd = exec.Command(SHUTDOWN_PROGRAM, "-h", "+1")
        default:
            cmd = exec.Command(SHUTDOWN_PROGRAM, "-h", "-P", "+1")
        }

        msg := SUCCESS_MSG
        err = cmd.Run()
        if err != nil {
            log.Printf("Error exec'ing: %v", err)
            msg = []byte("shutdownd|error|" + err.Error() + "\n")
        }

        _, err = sConn.WriteToUDP(msg, addr)
        if err != nil {
            log.Printf("Error writing: %v", err)
        }
    }
}
