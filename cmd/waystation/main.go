package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-waystation/internal/server";"github.com/stockyard-dev/stockyard-waystation/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9930"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./waystation-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("waystation: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Waystation — travel and trip planner\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("waystation: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
