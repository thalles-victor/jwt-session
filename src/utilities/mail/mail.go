package mail

import (
	"fmt"
	"jwt-session/src/utilities/config"
	"jwt-session/src/utilities/logger"

	gomail "gopkg.in/mail.v2"
)

func SendCreateAccount(name, email string) error {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "criação de conta na jwt-sesssion")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
			<meta charset="UTF-8" />
			<title>Bem-vindo!</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f7f7f7;
					color: #333;
					padding: 20px;
				}
				.container {
					background-color: #fff;
					border-radius: 8px;
					padding: 20px;
					box-shadow: 0 2px 5px rgba(0,0,0,0.1);
					max-width: 600px;
					margin: auto;
				}
				h2 {
					color: #0066cc;
				}
				p {
					font-size: 16px;
				}
				.footer {
					font-size: 12px;
					color: #777;
					margin-top: 20px;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Olá, %s 👋</h2>
				<p>Obrigado por criar uma conta na <strong>JWT-Session</strong>! 🎉</p>
				<p>Estamos muito felizes em tê-lo conosco.</p>
				<p>Em breve você poderá acessar todos os recursos da plataforma.</p>
				<div class="footer">
					<p>© 2025 JWT-Session. Todos os direitos reservados.</p>
				</div>
			</div>
		</body>
		</html>
	`, name)

	message.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(config.MAIL_HOST, config.MAIL_PORT, config.MAIL_USER, config.MAIL_PASS)

	logger.Info.Printf("send confirmation account to %s \n ", email)
	if err := dialer.DialAndSend(message); err != nil {
		logger.Error.Printf("error when send confirmation account email to %s", email)
		return err
	}
	logger.Info.Printf("Email sent to %s successfully!", email)

	return nil
}

func SendRecoveryRequestEmail(name, email, resetUrl string) error {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "recuperacao da conta na jwt-sesssion")

	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8" />
		<title>Recuperação de Conta</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f7f7f7;
				color: #333;
				padding: 20px;
			}
			.container {
				background-color: #fff;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0 2px 5px rgba(0,0,0,0.1);
				max-width: 600px;
				margin: auto;
			}
			h2 {
				color: #0066cc;
			}
			p {
				font-size: 16px;
			}
			.button {
				display: inline-block;
				padding: 12px 20px;
				margin: 20px 0;
				font-size: 16px;
				color: #fff;
				background-color: #0066cc;
				text-decoration: none;
				border-radius: 5px;
			}
			.footer {
				font-size: 12px;
				color: #777;
				margin-top: 20px;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Olá, %s 👋</h2>
			<p>Recebemos uma solicitação para redefinir sua senha da <strong>JWT-Session</strong>.</p>
			<p>Para continuar, clique no botão abaixo:</p>
			<a href="%s" class="button">Redefinir Senha</a>
			<p>Se você não solicitou a recuperação de senha, apenas ignore este email.</p>
			<div class="footer">
				<p>© 2025 JWT-Session. Todos os direitos reservados.</p>
			</div>
		</div>
	</body>
	</html>
`, name, resetUrl)

	message.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(config.MAIL_HOST, config.MAIL_PORT, config.MAIL_USER, config.MAIL_PASS)

	logger.Info.Printf("send confirmation account to %s \n ", email)
	if err := dialer.DialAndSend(message); err != nil {
		logger.Error.Printf("error when send confirmation account email to %s", email)
		return err
	}
	logger.Info.Printf("Email sent to %s successfully!", email)

	return nil
}
