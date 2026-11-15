package api

import (
 "bytes"
 "encoding/json"
 "image/png"
 "io"
 "net/http"
)
const maxLogoBytes = 5 << 20
const maxPixels = 16_000_000
type uploadResponse struct { Format string `json:"format"`; Width int `json:"width"`; Height int `json:"height"`; Bytes int `json:"bytes"` }
func uploadLogo(w http.ResponseWriter, r *http.Request) {
 r.Body = http.MaxBytesReader(w,r.Body,maxLogoBytes+4096)
 if err:=r.ParseMultipartForm(maxLogoBytes); err!=nil { http.Error(w,"invalid or oversized upload",413); return }
 defer r.MultipartForm.RemoveAll()
 files:=r.MultipartForm.File["logo"]
 if len(files)!=1 { http.Error(w,"exactly one logo required",400); return }
 if files[0].Size>maxLogoBytes { http.Error(w,"logo too large",413); return }
 f,err:=files[0].Open(); if err!=nil { http.Error(w,"cannot open file",400); return }; defer f.Close()
 data,err:=io.ReadAll(io.LimitReader(f,maxLogoBytes+1))
 if err!=nil || len(data)>maxLogoBytes { http.Error(w,"logo too large",413); return }
 if len(data)<8 || !bytes.Equal(data[:8],[]byte{137,80,78,71,13,10,26,10}) { http.Error(w,"PNG only in this milestone",415); return }
 cfg,err:=png.DecodeConfig(bytes.NewReader(data))
 if err!=nil || cfg.Width<1 || cfg.Height<1 || cfg.Width>maxPixels/cfg.Height { http.Error(w,"invalid dimensions",400); return }
 if _,err=png.Decode(bytes.NewReader(data)); err!=nil { http.Error(w,"corrupt PNG",400); return }
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(uploadResponse{"png",cfg.Width,cfg.Height,len(data)})
}
