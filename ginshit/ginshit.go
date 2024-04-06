package ginshit

import (
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type AuthTokenReq struct {
	KeyVaultUrl string `json:"key_vault_url"`
	SecretName  string `json:"secret_name"`
}

type AuthTokenResp struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

/*
https://learn.microsoft.com/zh-cn/azure/developer/go/azure-sdk-authentication-managed-identity?tabs=azure-cli#create-the-go-package
export AZURE_TENANT_ID=d68c51dc-5be6-40cb-87d5-f551a62e1049
export AZURE_CLIENT_ID=7dc60f4f-2426-4e12-9b9e-5136de73c01a
export AZURE_CLIENT_SECRET=zK1/lhC6v19J*lCJ2Dlm4TV^Wo#RX2xr

keyVaultURL = "https://shaprdcne2kvmgmtkv1.vault.azure.cn/"
secretName = "PSD86914-http"  // PSD86914-origin  PSD86630-http  PSD86630-origin
*/
func Main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(6 * time.Second)
		c.String(200, "Welcome to Gin!")
	})

	r.POST("/azure-auth-token", func(c *gin.Context) {
		var req AuthTokenReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, "parameter name cannot be empty")
			return
		}

		// keyVaultUrl, _ := c.Params.Get("key_vault_url")
		// secretName, _ := c.Params.Get("secret_name")
		log.Println("get azure token params", req.KeyVaultUrl, req.SecretName)

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Printf("failed to obtain a credential: %v\n", err)
			return
		}

		client, err := azsecrets.NewClient(req.KeyVaultUrl, cred, nil)
		if err != nil {
			log.Printf("failed to create a client: %v\n", err)
			return
		}

		resp, err := client.GetSecret(c.Request.Context(), req.SecretName, "", nil)
		if err != nil {
			log.Printf("failed to get a secret: %v\n", err)
			return
		}

		ret := AuthTokenResp{}
		// https://shaprdcne2kvmgmtkv1.vault.azure.cn/secrets/PSD86630-origin/c6c3291663e54da989b802765714743d
		if resp.ID != nil {
			ret.ID = string(*resp.ID)
		}

		// sr=c&si=PSD86630-racwl&sig=elhcIt6CZjbQ/X3yAbVoGtq7wugC68Ay2baCUH2X/FY=&sv=2020-06-12
		if resp.Value != nil {
			ret.Value = *resp.Value
		}

		c.JSON(http.StatusOK, ret)
	})

	// Use the pprof middleware
	pprof.Register(r)

	r.Run(":9999")
}
