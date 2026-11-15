package api
import (
 "bytes"
 "encoding/json"
 "image"
 "image/png"
 "mime/multipart"
 "net/http"
 "net/http/httptest"
 "testing"
)
func uploadRequest(t *testing.T, data []byte)*http.Request{
 t.Helper(); var b bytes.Buffer; m:=multipart.NewWriter(&b)
 f,err:=m.CreateFormFile("logo","logo.png"); if err!=nil { t.Fatal(err) }
 if _,err=f.Write(data);err!=nil { t.Fatal(err) }; if err=m.Close();err!=nil { t.Fatal(err) }
 r:=httptest.NewRequest(http.MethodPost,"/api/v1/logos/inspect",&b);r.Header.Set("Content-Type",m.FormDataContentType());return r
}
func TestPNG(t *testing.T){
 var b bytes.Buffer;if err:=png.Encode(&b,image.NewRGBA(image.Rect(0,0,2,3)));err!=nil{t.Fatal(err)}
 w:=httptest.NewRecorder();NewHandler().ServeHTTP(w,uploadRequest(t,b.Bytes()))
 if w.Code!=200 {t.Fatalf("status %d: %s",w.Code,w.Body.String())}
 var out uploadResponse;if err:=json.Unmarshal(w.Body.Bytes(),&out);err!=nil{t.Fatal(err)}
 if out.Width!=2||out.Height!=3||out.Format!="png"{t.Fatalf("unexpected: %+v",out)}
}
func TestUnsupported(t *testing.T){
 w:=httptest.NewRecorder();NewHandler().ServeHTTP(w,uploadRequest(t,[]byte("<svg/>")))
 if w.Code!=415{t.Fatalf("status %d",w.Code)}
}
func TestCorrupt(t *testing.T){
 w:=httptest.NewRecorder();NewHandler().ServeHTTP(w,uploadRequest(t,[]byte{137,80,78,71,13,10,26,10}))
 if w.Code!=400{t.Fatalf("status %d",w.Code)}
}
