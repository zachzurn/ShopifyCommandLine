# Go Shopify Tools

Go Shopify Tools is a Go project that provides tools for editing assets on shopify stores.

![go shopify tools](https://github.com/zachzurn/go_shopify_tools/raw/master/img/go_shopify_tools.png)

Functionality
  - Configure dev, test and staging stores [DONE]
  - Download entire theme [DONE]
  - Watch for file changes and upload to shopify [DONE]
  - Safe Deploy files to Shopify for a theme with incremental upload/deletes and deploy preview [NOT IMPLEMENTED]
  - Force Deploy all files to shopify with no confirmation. For automated use. [NOT IMPLEMENTED]
  - Deploy one store's theme to another store with 'source deploy target' command [NOT IMPLEMENTED] 
  - Copy products from one store to another [NOT IMPLEMENTED]
  - Import products using a JSON file [NOT IMPLEMENTED]
  - Encrypt and Decrypt shopify.json with password using 'settings lock' and 'settings unlock'. [NOT IMPLEMENTED]

### How to use

Directory structure

    Root Folder
    |--shopify.json
    |--shopify-dev (The name of this folder is set in the shopify.json stores section under the 'folder' field)
    |--shopify-staging (The name of this folder is set in the shopify.json stores section under the 'folder' field)
    |--shopify-production (The name of this folder is set in the shopify.json stores section under the 'folder' field)
    |--shopify.exe (I put the executable inside the folder, but you can run it from anywhere as long as you run the command from inside the root folder)

To configure your stores, you will need to create a shopify.json file that looks like the one below.

    {
    	"allowed_extensions" : [".txt"],
    	"stores" : {
    		"dev" : {
    			"name" : "Shopify Dev",
    			"folder" : "shopify-dev",
    			"api_key": "STORE_API_KEY_GOES_HERE",
    			"api_password": "STORE_API_PASSWORD_GOES_HERE",
    			"shop": "STORE_URL_GOES_HERE.myshopify.com",
    			"theme_id" : 00000000
    		},
    		"staging" : {
    			"name" : "Shopify Staging",
    			"folder" : "shopify-staging",
    			"api_key": "STORE_API_KEY_GOES_HERE",
    			"api_password": "STORE_API_PASSWORD_GOES_HERE",
    			"shop": "STORE_URL_GOES_HERE.myshopify.com",
    			"theme_id" : 00000000
    		},
    		"production" : {
    			"name" : "Shopify Production",
    			"folder" : "shopify-production",
    			"api_key": "STORE_API_KEY_GOES_HERE",
    			"api_password": "STORE_API_PASSWORD_GOES_HERE",
    			"shop": "STORE_URL_GOES_HERE.myshopify.com",
    			"theme_id" : 00000000
    		},
    	}
    }
    
You can create as many stores as you want and name them however you want. Once you have your settings in place you can run the following commands.

**Watch**

Watches the folder specified in the "folder" field for changes and uploads/deletes changes.

    watch [store] //Windows Example ./shopify.exe watch dev
    
**Download**

Downloads the entire theme to the folder specified in the "folder" field.
    
    download [store] //Windows Example ./shopify.exe download production

### Development

Want to contribute?

- See Functionality above and help with [NOT IMPLEMENTED] features
- Test on OSX and Linux (Tested on Windows so far) and report/fix issues