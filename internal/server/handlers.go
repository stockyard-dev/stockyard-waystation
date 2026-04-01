package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-waystation/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Trip{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var t store.Trip;json.NewDecoder(r.Body).Decode(&t);if t.Name==""{writeError(w,400,"name required");return};if t.Status==""{t.Status="planned"};s.db.Create(&t);writeJSON(w,201,t)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListLegs(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListLegs(id);if list==nil{list=[]store.Leg{}};writeJSON(w,200,list)}
func(s *Server)handleAddLeg(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var l store.Leg;json.NewDecoder(r.Body).Decode(&l);l.TripID=id;if l.Description==""{writeError(w,400,"description required");return};s.db.AddLeg(&l);writeJSON(w,201,l)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
