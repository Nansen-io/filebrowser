# About

This markdown file outlines the plan of forking the file browser project 

https://github.com/gtsteffaniak/filebrowser

to integrate with the ChainFS API (I am the developer behind this project): 

- DEV: https://nansendev.azurewebsites.net/swagger/v1/swagger.json 
- UAT: https://nansenuat.azurewebsites.net/swagger/v1/swagger.json
- PROD: https://nansenprod.azurewebsites.net/swagger/v1/swagger.json

Full source code is available to reference in the local path:

C:\git\azure-blockchain-workbench-app\NasenAPI

# Goals

## Authentication (Priority 1)

Update filebrowser to use Chainfs authentication which is Azure AD B2C (note not entra).

the chainfs API has the endpoint /api/NansenFile/LoginURL that gets the login URL for the current environment.

```bash
curl -X 'GET' \
  'https://nansendev.azurewebsites.net/api/NansenFile/LoginURL' \
  -H 'accept: text/plain'
```  

Auth URL for DEV: https://NansenDEV2.b2clogin.com/NansenDEV2.onmicrosoft.com/B2C_1_signupsignin1/oauth2/v2.0/authorize?client_id=ae8e4cce-f313-459b-b86b-2fa59b4f1cb8&redirect_uri=https%3a%2f%2fjwt.ms%2f&response_type=token&scope=openid+profile+offline_access+https%3a%2f%2fNansenDEV2.onmicrosoft.com%2ftasks-api%2faccess_as_user&response_mode=fragment

Logout URL can be retrieved from the API using the endpoint /api/NansenFile/LogoutURL

```bash
curl -X 'GET' \
  'https://nansendev.azurewebsites.net/api/NansenFile/LogoutURL' \
  -H 'accept: text/plain'
```

Logout URL (DEV): https://nansendev.azurewebsites.net:/api/User/Logout

## Right Click Menu Addition (Priority 2)

Add right click menu to Filebrowser that calls chainfs API endpoint: /api/NansenFile/FileCreate

The chainfs API has a debug section that provides a reference implementation of encoding a file submission for the FileCreate endpoint:

/api/Debug/FileEncode

Expected endpoints needed to sync filebrowser with chainfs API:

- /api/NansenFile/DirCreate
- /api/NansenFile/DirGetInfo
- /api/NansenFile/DirRename
- /api/NansenFile/DirSubDirs
- /api/NansenFile/FileCreate
- /api/NansenFile/FileExists
- /api/NansenFile/FileGetDetails
- /api/NansenFile/FileNewest
- /api/NansenFile/FileUpdate
- /api/NansenFile/GetFileSimpleInfo
- /api/NansenFile/ListOfDirectories
- /api/NansenFile/SetTags

File right click menu behaviour:

- Store on Chainfs
- Update file on Chainfs (when file has changed)

as we are going for a MVP our goal is to simply take a file from filebrowser and store it on chainfs, but we should design with the future in mind.

If FileBrowser uses UUIDs to track its files, we should that as the filename on chainfs. this part will need planning.

Tracking file updates will need consideration as well. 

Chainfs uses the concept of a **genesisGuid** which is the first revision in a files history. this can be found for any file and or revision already on chainfs by using the end point /api/NansenFile/GetFileSimpleInfo

then for submitting an update, you would use /api/NansenFile/FileNewest to get the latest GUID for the genesisGuid (or any file GUID in the history) to get the guidValue of the latest revision of the file stored on chainfs.

it is not neccesary to sync all changes from filebrowser to chainfs, but users should have an indicator if the latest revision of their file on filebrowser is not stored as an update on chainfs.

## Hosting  (Priority 3)

Host Fork on azure. 

Three instances will be created: DEV, UAT & PROD that each point to the relevant ChainFS environment.