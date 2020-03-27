package service

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/AI-Research-HIT/2019-nCoV-Service/config"
	"github.com/ender-wan/ewlog"
	"github.com/gorilla/mux"

	mw "github.com/AI-Research-HIT/2019-nCoV-Service/middleware"

	"github.com/AI-Research-HIT/2019-nCoV-Service/auth"
)

func StartService(ctx context.Context) {
	auth.InitJwt()

	r := mux.NewRouter()
	//AttachProfiler(r)

	r.HandleFunc("/login", LoginHandler).Methods("POST")

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/model-cal", ModelCalculateHandler).Methods("POST")
	apiRouter.HandleFunc("/mdata", MDataHandler).Methods("GET")
	apiRouter.HandleFunc("/latest", LatestDataHandler).Methods("GET")
	apiRouter.HandleFunc("/mlist", MlistHandler).Methods("GET")
	apiRouter.HandleFunc("/chinamobile-data", ChinaMobileDataHandler).Methods("GET")
	apiRouter.HandleFunc("/baidu-incity", BaiduInCityHandler).Methods("GET")
	apiRouter.HandleFunc("/province-data", AllProvinceDataHandler).Methods("GET")
	apiRouter.HandleFunc("/simulate", MonteCarloSimulationHandler).Methods("POST")
	apiRouter.HandleFunc("/allRegionName", GetAllProvinceOrCountryName).Methods("GET")
	apiRouter.HandleFunc("/regionDailyData", GetProvinceDailyData).Methods("GET")
	apiRouter.HandleFunc("/countryCompare", CompareCountryTrending).Methods("GET")
	apiRouter.HandleFunc("/m-simulate", CompareWithDiffMHandler).Methods("POST")

	r.Use(mw.CorsMw)
	//r.Use(mw.JwtAuthMw)

	go func() {
		err := http.ListenAndServe(config.Config.ServerAddr, r)
		if err != nil {
			ewlog.Fatal(err)
		}
	}()
}

func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}
