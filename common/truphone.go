package common

import(
    "net/http"
    "net/url"
    "strconv"
    "errors"
)

type Truphone struct {

}

func (self* Truphone) SendSms(to string, from string, msg string) error{
    path := "https://api.tp.mu/dataart-out.php"
    resp, err := http.PostForm(path,
        url.Values{
            "from":    {from},
            "to":      {to},
            "message": {msg}})

    if err != nil {
        return err
    }

    if resp.StatusCode != 200 {
        return errors.New("Resp code is not 200 for " + path + "; StatusCode = " + strconv.Itoa(resp.StatusCode))
    }
    return nil
}