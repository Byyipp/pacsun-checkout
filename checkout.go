package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"strings"

	"github.com/bwmarrin/discordgo"
)

const token string = "ODQ2NDUwMDg3MzkxNTI2OTUy.YKvsEw.24njjK8Q4tCJ6KF7X8RtRaeT6bY"

var BotID string

type prod struct {
	Products []vars `json:"Products"`
}

type vars struct {
	Size            string `json:"Size"`
	Color           string `json:"Color"`
	ProductURL      string `json:"ProductURL"`
	OnlineInventory int    `json:"OnlineInventory"`
	Availability    string `json:"Availability"`
	ImageURL        string `json:"ImageURL"`
	Title           string `json:"Title"`
	ProductID       string `json:"ProductID"`
}

type mapper struct {
	color    string
	size     string
	URL      string
	imageURL string
	title    string
	ID       string
}

type atcData struct {
	cartAction string
	Quantity   string
	pid        string
}

type profileData struct {
	Profile_securekey                                         string `json:"dwfrm_profile_securekey"`
	Billing_billingAddress_addressFields_email_emailAddress   string `json:"dwfrm_billing_billingAddress_addressFields_email_emailAddress"`
	Billing_billingAddress_addressFields_phone                string `json:"dwfrm_billing_billingAddress_addressFields_phone"`
	Singleshipping_shippingAddress_optInEmail                 string `json:"dwfrm_singleshipping_shippingAddress_optInEmail"`
	Singleshipping_shippingAddress_alternateFirstName         string `json:"dwfrm_singleshipping_shippingAddress_alternateFirstName"`
	Singleshipping_shippingAddress_alternateLastName          string `json:"dwfrm_singleshipping_shippingAddress_alternateLastName"`
	Singleshipping_securekey                                  string `json:""`
	Singleshipping_shippingAddress_addressFields_firstName    string `json:"dwfrm_singleshipping_securekey"`
	Singleshipping_shippingAddress_addressFields_lastName     string `json:"dwfrm_singleshipping_shippingAddress_addressFields_lastName"`
	Singleshipping_shippingAddress_addressFields_address1     string `json:"dwfrm_singleshipping_shippingAddress_addressFields_address1"`
	Singleshipping_shippingAddress_addressFields_address2     string `json:"dwfrm_singleshipping_shippingAddress_addressFields_address2"`
	Singleshipping_shippingAddress_addressFields_city         string `json:"dwfrm_singleshipping_shippingAddress_addressFields_city"`
	Singleshipping_shippingAddress_addressFields_states_state string `json:"dwfrm_singleshipping_shippingAddress_addressFields_states_state"`
	Singleshipping_shippingAddress_addressFields_country      string `json:"dwfrm_singleshipping_shippingAddress_addressFields_country"`
	Singleshipping_shippingAddress_addressFields_postal       string `json:"dwfrm_singleshipping_shippingAddress_addressFields_postal"`
	Singleshipping_originID                                   string `json:"dwfrm_singleshipping_originID"`
	Billing_paymentMethods_selectedPaymentMethodID            string `json:"dwfrm_billing_paymentMethods_selectedPaymentMethodID"`
	Billing_paymentMethods_creditCard_number                  string `json:"dwfrm_billing_paymentMethods_creditCard_number"`
	Billing_paymentMethods_creditCard_owner                   string `json:"dwfrm_billing_paymentMethods_creditCard_owner"`
	Billing_paymentMethods_creditCard_type                    string `json:"dwfrm_billing_paymentMethods_creditCard_type"`
	ExpDate                                                   string `json:"expDate"`
	Billing_paymentMethods_creditCard_expiration_month        string `json:"dwfrm_billing_paymentMethods_creditCard_expiration_month"`
	Billing_paymentMethods_creditCard_expiration_year         string `json:"dwfrm_billing_paymentMethods_creditCard_expiration_year"`
	Billing_paymentMethods_creditCard_cvn                     string `json:"dwfrm_billing_paymentMethods_creditCard_cvn"`
	Billing_save                                              string `json:"dwfrm_billing_save"`
	Billing_securekey                                         string `json:"dwfrm_billing_securekey"`
	BcriptionCode                                             string `json:"ltkSubscriptionCode"`
	Billing_billingAddress_addressFields_firstName            string `json:"dwfrm_billing_billingAddress_addressFields_firstName"`
	Billing_billingAddress_addressFields_lastName             string `json:"dwfrm_billing_billingAddress_addressFields_lastName"`
	Billing_billingAddress_addressFields_address1             string `json:"dwfrm_billing_billingAddress_addressFields_address1"`
	Billing_billingAddress_addressFields_address2             string `json:"dwfrm_billing_billingAddress_addressFields_address2"`
	Billing_billingAddress_addressFields_city                 string `json:"dwfrm_billing_billingAddress_addressFields_city"`
	Billing_billingAddress_addressFields_states_state         string `json:"dwfrm_billing_billingAddress_addressFields_states_state"`
	Billing_billingAddress_addressFields_postal               string `json:"dwfrm_billing_billingAddress_addressFields_postal"`
	Billing_billingAddress_addressFields_country              string `json:"dwfrm_billing_billingAddress_addressFields_country"`
	Csrf_token                                                string `json:"csrf_token"`
	ShippingID                                                string `json:"shippingID"`
}

func createwebhook(list [][]mapper) *discordgo.WebhookParams {
	var thumbnail string
	test := false
	for i := range list {
		for k := range list[i] {
			if strings.Compare(list[i][k].imageURL, "") != 0 {
				thumbnail = list[i][k].imageURL
				test = true
				break
			}
		}
		if test {
			break
		}
	}
	var titty string
	test = false
	for i := range list {
		for k := range list[i] {
			if strings.Compare(list[i][k].title, "") != 0 {
				titty = list[i][k].title
				test = true
				break
			}
		}
		if test {
			break
		}
	}
	var fields []*discordgo.MessageEmbedField
	for i := range list {
		temp := &discordgo.MessageEmbedField{
			Name:   "",
			Value:  "",
			Inline: true,
		}
		temp.Name = list[i][0].color

		for k := range list[i] {
			temp.Value = temp.Value + "[" + list[i][k].size + "](" + list[i][k].URL + ")\n"
		}
		fields = append(fields, temp)

	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Pacsun.com",
			URL:  "https://www.pacsun.com/",
		},
		Title: "Product Loaded",
		URL:   list[0][0].URL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnail,
		},
		Description: titty,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Nuggie AIO - " + time.Now().String(),
			IconURL: "https://cdn.discordapp.com/emojis/839655842705571862.png?v=1",
		},
		Color: 16549376,
	}
	var id []*discordgo.MessageEmbed
	id = append(id, embed)
	webhook := discordgo.WebhookParams{
		Username: "Nuggie AIO",
		Embeds:   id,
	}

	return &webhook

}

func main() {
	rand.Seed(time.Now().UnixNano())
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	_ = dg.Open()

	variant := "0172499380015" //PRODUCT VARIANT INSERT HERE

	var params prod

	for { //monitoring method
		c := http.Client{}
		req, err := http.NewRequest("GET", "https://cloudservices.predictspring.com/search/v1/variantgroup/"+variant, strings.NewReader(""))
		req.Header = http.Header{
			"Host":                         {"cloudservices.predictspring.com"},
			"Accept":                       {"*/*"},
			"firstLogin":                   {"true"},
			"PredictSpring-Region":         {"US"},
			"PredictSpring-API-Key":        {"BEDpgdZPPnEC4RWKrwHGJDNN"},
			"PredictSpring-Locale":         {"en_US"},
			"PredictSpring-DeviceName":     {"E"},
			"PredictSpring-MerchantID":     {"PACSUN453239"},
			"Accept-Language":              {"en-us"},
			"PredictSpring-Attribution":    {"%7B%0A%20%20%22appVersion%22%20:%20%225.2.0%22,%0A%20%20%22appBuildVersion%22%20:%20%2215541%22,%0A%20%20%22deviceOS%22%20:%20%22IOS%22%0A%7D"},
			"User-Agent":                   {"Pacsun/15541 CFNetwork/1121.2.2 Darwin/19.3.0"},
			"Connection":                   {"keep-alive"},
			"PredictSpring-InstallationID": {"b3752662-225e-4b24-a4b9-4f5ebc591674"},
		}
		if err != nil {
			println(err)
		}

		resp, err := c.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err)
		}

		//println("body:\n" + string(body))
		//println(len(string(body))) //check if stock loaded?
		if len(string(body)) < 100 {
			println("not loaded")
			time.Sleep(time.Second * 7)
			continue
		}

		err = json.Unmarshal([]byte(string(body)), &params)
		if err != nil {
			log.Fatal(err)
		}

		break

	}
	slice := make([]mapper, 0)
	for i := range params.Products {
		temp := mapper{
			ID: params.Products[i].ProductID,
		}
		slice = append(slice, temp)
	}

	num := rand.Intn(len(slice))
	println(slice[num].ID)

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Print("ERROR CREATING COOKIE JAR :(")
	}

	c := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 30,
	}

	data := url.Values{}
	data.Set("cartAction", "add")
	data.Set("Quantity", "1")
	data.Set("pid", slice[num].ID)

	r, err := http.NewRequest("POST", "https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/Cart-AddProduct?format=ajax", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	post, err := c.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(post.Status)

	defer post.Body.Close()

	// poststr, err := ioutil.ReadAll(post.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// println(string(poststr))

	//display response^^^^

	//SHIPPING
	data = url.Values{}
	data.Set("address1", "101 test st")
	data.Set("city", "adrian")
	data.Set("countryCode", "US")
	data.Set("postalCode", "49221-4462")
	data.Set("stateCode", "MI")
	data.Set("responseObject", "getShippingMethods")
	data.Set("ajax", "true")

	r, err = http.NewRequest("POST", "https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COCheckout-GetShippingMethods?format=ajax", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	post, err = c.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(post.Status)

	defer post.Body.Close()

	//VERIFY
	data = url.Values{}
	data.Set("firstName", "baby")
	data.Set("lastName", "Yip")
	data.Set("address1", "1700 Harold St")
	data.Set("address2", "")
	data.Set("city", "Adrian")
	data.Set("state", "MI")
	data.Set("zip", "49221-4462")
	data.Set("country", "USA")

	r, err = http.NewRequest("POST", "https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/VerifyOneStepAddress-ShippingAddress", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	post, err = c.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(post.Status)

	defer post.Body.Close()

	//CHECKOUT
	// checkout := profileData{
	// 	Profile_securekey: "1416149199",
	// 	Billing_billingAddress_addressFields_email_emailAddress:   "tets@gmail.com",
	// 	Billing_billingAddress_addressFields_phone:                "201 5555555",
	// 	Singleshipping_shippingAddress_optInEmail:                 "true",
	// 	Singleshipping_securekey:                                  "385646491",
	// 	Singleshipping_shippingAddress_addressFields_firstName:    "baby",
	// 	Singleshipping_shippingAddress_addressFields_lastName:     "yip",
	// 	Singleshipping_shippingAddress_addressFields_address1:     "1700 Harold St",
	// 	Singleshipping_shippingAddress_addressFields_city:         "Adrian",
	// 	Singleshipping_shippingAddress_addressFields_states_state: "MI",
	// 	Singleshipping_shippingAddress_addressFields_country:      "US",
	// 	Singleshipping_shippingAddress_addressFields_postal:       "49221-4462",
	// 	Singleshipping_originID:                                   "DSK",
	// 	Billing_paymentMethods_selectedPaymentMethodID:            "CREDIT_CARD",
	// 	Billing_paymentMethods_creditCard_number:                  "4242424242424242",
	// 	Billing_paymentMethods_creditCard_owner:                   "harold dr",
	// 	Billing_paymentMethods_creditCard_type:                    "Visa",
	// 	ExpDate:                                                   "01/25",
	// 	Billing_paymentMethods_creditCard_expiration_month:        "1",
	// 	Billing_paymentMethods_creditCard_expiration_year:         "2025",
	// 	Billing_paymentMethods_creditCard_cvn:                     "666",
	// 	Billing_save:                                              "true",
	// 	Billing_securekey:                                         "1509941636",
	// 	BcriptionCode:                                             "checkoutbilling",
	// 	Billing_billingAddress_addressFields_firstName:            "baby",
	// 	Billing_billingAddress_addressFields_lastName:             "yip",
	// 	Billing_billingAddress_addressFields_address1:             "1700 Harold St",
	// 	Billing_billingAddress_addressFields_city:                 "Adrian",
	// 	Billing_billingAddress_addressFields_states_state:         "MI",
	// 	Billing_billingAddress_addressFields_postal:               "49221-4462",
	// 	Billing_billingAddress_addressFields_country:              "US",
	// 	Csrf_token: "ZrB2QunW_2YV3SKNX0Thxq5rDHzDUBvAJkvy8Dmskfy5K1MCwuR6pWP_GboI7zFquCASTTeHiw1fCwMI5Xaj7NfvdAAoMrcpvt9dOhsVlv5vEhSu-KmPmAxIX-u-ZhtWCHNqCa4oBuYNbERPMPZ3dLRrunW32SuOtXJuuPmqvTYDWoKpULA=",
	// 	ShippingID: "SP",
	// }

	// jsonCheckout, _ := json.Marshal(checkout)
	// postcheckout, err := c.Post("https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COCheckout-OrderSubmit", "application/json", bytes.NewBuffer(jsonCheckout))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// //checkoutbody, _ := ioutil.ReadAll(postcheckout.Body)
	// postcheckout.Body.Close()
	// println(postcheckout.Status)
	// //println(string(checkoutbody))

	// // f, err := os.OpenFile("body.txt", os.O_APPEND|os.O_WRONLY, 0600)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // defer f.Close()
	// // if _, err = f.WriteString(string(checkoutbody)); err != nil {
	// // 	log.Fatal(err)
	// // }

	// GET SECUREKEYS

	url, err := url.Parse("https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COSummary-Submit")
	if err != nil {
		fmt.Println(err)
	}

	a := jar.Cookies(url)
	// for i := range a {
	// 	println(a[i].String())
	// }

	var cookiestring string
	for i := range a {
		if i == len(a)-1 {
			cookiestring = cookiestring + a[i].String()
			break
		}
		cookiestring = cookiestring + a[i].String() + "; "
	}
	//println(cookiestring)

	r, err = http.NewRequest("GET", "https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COSummary-Submit", strings.NewReader(""))
	r.Header = http.Header{
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"accept-encoding":           {"gzip, deflate, br"},
		"accept-language":           {"en-US,en;q=0.9"},
		"cookie":                    {cookiestring},
		"dnt":                       {"1"},
		"referer":                   {"https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COSummary-Submit"},
		"sec-ch-ua":                 {`" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"same-origin"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"},
	}
	if err != nil {
		log.Fatal(err)
	}

	post, err = c.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(post.Status)

	defer post.Body.Close()
	// poststr, err := ioutil.ReadAll(post.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// println(string(poststr))

	// r, err = http.NewRequest("POST", "https://www.pacsun.com/on/demandware.store/Sites-pacsun-Site/default/COCheckout-OrderSubmit", strings.NewReader(""))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// post, err = c.Do(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(post.Status)

	// defer post.Body.Close()
	// poststr, err := ioutil.ReadAll(post.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// println(string(poststr))

}
