package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/imhasandl/crypto-management/database"
	"github.com/imhasandl/crypto-management/utils"
)

type Config struct {
	db     *database.DB
	runner map[string]chan struct{}
}

func NewConfig(db *database.DB) *Config {
	return &Config{
		db:     db,
		runner: make(map[string]chan struct{}),
	}
}

// AddCurrency godoc
// @Summary Добавить монету для отслеживания
// @Description Запускает runner для монеты и начинает получать цену
// @Tags currency
// @Accept json
// @Produce json
// @Param coin body handlers.CoinRequest true "CoinGecko ID монеты"
// @Success 201 {object} map[string]string
// @Failure 400,409,500 {object} map[string]string
// @Router /currency/add [post]
func (cfg *Config) AddCurrency(w http.ResponseWriter, req *http.Request) {
	var request CoinRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "can't decode request", err)
		return
	}

	_, exists := cfg.runner[request.Coin]
	if exists {
		utils.RespondWithError(w, http.StatusConflict, "already running crypto", nil)
		return
	}

	stop := make(chan struct{})
	cfg.runner[request.Coin] = stop
	go cfg.startRunner(request.Coin, stop)

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Currency added and runner started",
		"coin":    request.Coin,
	})
}

func (cfg *Config) startRunner(coin string, stop chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			price, err := fetchPrice(coin)
			if err != nil {
				fmt.Printf("failed to fetch price for %s: %v\n", coin, err)
				continue
			}
			err = cfg.db.SaveCoinPrice(coin, price, time.Now())
			if err != nil {
				fmt.Printf("failed to save price for %s: %v\n", coin, err)
			}

			log.Printf("%v added, price %v$", coin, price)
		case <-stop:
			return
		}
	}
}

func fetchPrice(coin string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coin))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	fmt.Println("CoinGecko response:", string(body))

	var result map[string]map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	coinData, ok := result[coin]
	if !ok {
		return 0, fmt.Errorf("coin '%s' not found in API", coin)
	}

	val, ok := coinData["usd"]
	if !ok {
		return 0, fmt.Errorf("USD price not found for coin '%s'", coin)
	}

	switch v := val.(type) {
	case float64:
		return int(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse string price: %v", err)
		}
		return int(f), nil
	default:
		return 0, fmt.Errorf("unexpected type for price: %T", v)
	}
}

// RemoveCurrency godoc
// @Summary Удалить монету из отслеживания
// @Description Останавливает runner для монеты
// @Tags currency
// @Accept json
// @Produce json
// @Param coin body handlers.CoinRequest true "CoinGecko ID монеты"
// @Success 200 {object} map[string]string
// @Failure 400,409,500 {object} map[string]string
// @Router /currency/remove [post]
func (cfg *Config) RemoveCurrency(w http.ResponseWriter, req *http.Request) {
	var request CoinRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "can't decode request", err)
		return
	}

	stop, exists := cfg.runner[request.Coin]
	if !exists {
		utils.RespondWithError(w, http.StatusConflict, "coin doesn't exists", err)
		return
	}

	close(stop)
	delete(cfg.runner, request.Coin)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Currency removed and runner stopped",
		"coin":    request.Coin,
	})
}

// GetCurrencyPrice godoc
// @Summary Получить цену монеты на определённое время
// @Description Возвращает цену монеты, ближайшую к указанному времени
// @Tags currency
// @Accept json
// @Produce json
// @Param coin body handlers.CoinPriceRequest true "CoinGecko ID монеты и timestamp"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]string
// @Router /currency/price [post]
func (cfg *Config) GetCurrencyPrice(w http.ResponseWriter, req *http.Request) {
	var request CoinPriceRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "can't decode request", err)
		return
	}

	price, ts, err := cfg.db.GetNearestPrice(request.Coin, time.Unix(request.Timestamp, 0))
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "price not found", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"coin":      request.Coin,
		"price":     price,
		"timestamp": ts.Unix(),
	})
}
