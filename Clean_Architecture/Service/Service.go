package Service

/*
const ApiKey string = "22d3fa86b8743e497b32195cbc690abc06b42436"
const SecretKey string = "adf07bdd63b240ae60087efd2e72269b9c65cc91"

func Register(User io.ReadCloser) (s.NewUserResponse, error) {

	var regData s.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(User)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось прочитать запрос")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось дессериализировать JSON")
	}
	tokenAuth, err := TokenGenerate(regData.Email, regData.Password)
	if err != nil {
		return s.NewUserResponse{}, err
	}
	// Устанавливаем токен и добавляем пользователя в БД
	regData.TokenString.Token = tokenAuth

	err = Repository.Create(regData)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось добавить нового пользователя в БД")
	}

	return regData, nil

}
func Login(Usertoken string, User io.ReadCloser) (s.NewUserResponse, error) {

	var regData s.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(User)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось прочитать запрос")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось дессериализировать JSON")
	}
	// Проверяем логин, пароль и токен
	Email, Password, Token := UserInfo_Checker(regData.Email, regData.Password, Usertoken)
	if !Email {
		return s.NewUserResponse{}, errors.New("не верный логин")
	}
	if !Password {
		return s.NewUserResponse{}, errors.New("не верный пароль")
	}
	if !Token {
		freshToken := RefreshToken(regData.Email, regData.Password)
		regData.TokenString.Token = freshToken
		return s.NewUserResponse{}, errors.New("вы успешно вышли из сервиса")
	}

	return regData, nil
}

func TokenGenerate(email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": password})
	if err != nil {
		return "regData", errors.New("ошибка генерации токена")
	}
	return tokenString, nil
}
func RefreshToken(email, pawwsord string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": pawwsord})
	if err != nil {
		log.Println(err)
	}
	err = Repository.RefreshToken(email, pawwsord, tokenString)
	if err != nil {
		log.Println(err)
	}
	return tokenString
}

func UserInfo_Checker(email, password, token string) (bool, bool, bool) {
	var Email, Password, Token bool
	Email = Repository.CheckEmail(email)
	Password = Repository.CheckPassword(password)
	Token = Repository.CheckToken(token)
	return Email, Password, Token
}

func Search(UserRequest io.ReadCloser) (s.RequestQuery, error) {

	bodyJSON, err := ioutil.ReadAll(UserRequest)
	if err != nil {
		return s.RequestQuery{}, errors.New("не удалось прочитать запрос")
	}

	var SearchResp s.RequestUser
	var requestQuery s.RequestQuery

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		return s.RequestQuery{}, errors.New("не удалось дессериализировать JSON")
	}

	resp, err := HandleWorker(SearchResp)
	if err != nil {
		return s.RequestQuery{}, err
	}
	requestQuery.Query = fmt.Sprintf("Широта: %s, Долгота: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return requestQuery, nil
}

func Address(UserRequest io.ReadCloser) (s.RequestQuery, error) {

	bodyJSON, err := ioutil.ReadAll(UserRequest)
	if err != nil {
		return s.RequestQuery{}, errors.New("не удалось прочитать запрос")
	}

	var SearchResp s.RequestUser
	var requestQuery s.RequestQuery

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		return s.RequestQuery{}, errors.New("не удалось дессериализировать JSON")
	}

	resp, err := HandleWorker(SearchResp)
	if err != nil {
		return s.RequestQuery{}, err
	}
	requestQuery.Query = fmt.Sprintf("Отформатированный адресс %s", resp.Addres)

	return requestQuery, nil
}

func HandleWorker(Qwery s.RequestUser) (s.RequestAddress, error) {
	var requestQuery s.RequestAddress
	cache, err := Repository.CacheChecker(Qwery)
	if err != nil {
		log.Println(err)
	}
	if cache {
		requestQuery.Addres = cache.Addres
		requestQuery.RequestSearch.Lat = cache.RequestSearch.Lat
		requestQuery.RequestSearch.Lng = cache.RequestSearch.Lng
		return requestQuery, nil
	}

	geocodeResponse, err := Geocode(Qwery.RequestQuery)
	if err != nil {
		return s.RequestAddress{}, errors.New("ошибка в работе dadata")
	}
	for _, v := range geocodeResponse {
		requestQuery.RequestSearch.Lat = v.GeoLat
		requestQuery.RequestSearch.Lng = v.GeoLon
		requestQuery.Addres = v.Result
	}
	err = Repository.Select(Qwery)
	if err != nil {
		return s.RequestAddress{}, errors.New("ошибка запроса Select")
	}
	return requestQuery, nil
}

func Geocode(Querys s.RequestQuery) ([]*model.Address, error) {

	creds := client.Credentials{
		ApiKeyValue:    ApiKey,
		SecretKeyValue: SecretKey,
	}

	api := dadata.NewCleanApi(client.WithCredentialProvider(&creds))

	result, err := api.Address(context.Background(), Querys.Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil

}

*/
