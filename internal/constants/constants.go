package constants

const (
	// Fichiers de configuration
	ConfigFile  = "config.json"
	CookiesFile = "cookies.json"

	// Version de l'application
	AppVersion   = "other"
	DebugEnabled = false

	// Messages utilisateur
	CaptchaPromptMessage   = "Please solve CAPTCHA at the following address : %s\n"
	CaptchaTokenPrompt     = "Enter CAPTCHA : "
	CaptchaRequiredMessage = "CAPTCHA verification required"
	TOTPAutoMessage        = "Automatic code generation TOTP..."
	TOTPPrompt             = "Enter TOTP : "

	// Méthode de vérification
	VerificationMethodCaptcha = "captcha"

	// Code d'erreur CAPTCHA
	CaptchaErrorCode = 9001
)
