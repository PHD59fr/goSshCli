# go sshcli

Please Change the MAIN:

    sshCli := &SSH{
	    Host :  **YOURHOSTNAME**,
	    User :  **YOURUSERNAME**,
	    Port :  22,
	    Credential: homeDir + "/.ssh/id_rsa",
    }
    sshCli.Connect(CERT_PUBLIC_KEY_FILE) 

You can use CERT_PUBLIC_KEY_FILE or CERT_PASSWORD but if you choice password, please change the Credential value
