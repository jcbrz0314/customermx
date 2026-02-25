package mail

import "fmt"

func emailLayout(title, content, buttonLabel, buttonURL, buttonColor string) string {
	buttonHTML := ""
	if buttonLabel != "" && buttonURL != "" {
		buttonHTML = fmt.Sprintf(`
                        <table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin: 28px auto 0;">
                            <tr>
                                <td style="border-radius: 6px; background-color: %s;">
                                    <a href="%s" target="_blank" style="display: inline-block; padding: 14px 32px; font-size: 14px; font-family: Arial, sans-serif; font-weight: 600; color: #ffffff; text-decoration: none; border-radius: 6px;">
                                        %s
                                    </a>
                                </td>
                            </tr>
                        </table>`, buttonColor, buttonURL, buttonLabel)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
</head>
<body style="margin: 0; padding: 0; background: linear-gradient(135deg, #1a237e 0%%, #283593 50%%, #3949ab 100%%); font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif;">
    <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background: linear-gradient(135deg, #1a237e 0%%, #283593 50%%, #3949ab 100%%); padding: 40px 16px;">
        <tr>
            <td align="center">
                <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="560" style="max-width: 560px; width: 100%%;">

                    <!-- Brand header -->
                    <tr>
                        <td align="center" style="padding-bottom: 28px;">
                            <p style="margin: 0; font-size: 30px; font-weight: 900; letter-spacing: 1px; color: #ffffff; font-family: Arial, sans-serif;">
                                Customer<span style="color: #90caf9;">MX</span>
                            </p>
                            <p style="margin: 5px 0 0; font-size: 10px; letter-spacing: 4px; text-transform: uppercase; color: rgba(255,255,255,0.6); font-family: Arial, sans-serif;">
                                Gestión de Eventos Automotrices
                            </p>
                        </td>
                    </tr>

                    <!-- Card -->
                    <tr>
                        <td style="background-color: #ffffff; border-radius: 12px; box-shadow: 0 10px 40px rgba(0,0,0,0.25); padding: 40px 36px;">

                            <h1 style="margin: 0 0 8px; font-size: 22px; font-weight: 700; color: #1a237e; text-align: center;">
                                %s
                            </h1>

                            <hr style="border: none; border-top: 2px solid #3949ab; width: 48px; margin: 16px auto 28px;" />

                            %s
                            %s

                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="padding: 24px 16px 0; text-align: center;">
                            <p style="margin: 0; font-size: 12px; color: rgba(255,255,255,0.6); line-height: 1.5;">
                                CustomerMX &middot; Gestión de Eventos Automotrices
                            </p>
                            <p style="margin: 6px 0 0; font-size: 11px; color: rgba(255,255,255,0.4);">
                                Este es un correo automático, por favor no respondas a este mensaje.
                            </p>
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, title, title, content, buttonHTML)
}

func infoBox(rows string) string {
	return fmt.Sprintf(`
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color: #F5F7FA; border-radius: 8px; margin: 20px 0; padding: 0;">
                                <tr>
                                    <td style="padding: 16px 20px;">
                                        %s
                                    </td>
                                </tr>
                            </table>`, rows)
}

func infoRow(label, value string) string {
	return fmt.Sprintf(`<p style="margin: 4px 0; font-size: 14px; color: #435561;"><strong>%s</strong> %s</p>`, label, value)
}

func paragraph(text string) string {
	return fmt.Sprintf(`<p style="margin: 0 0 16px; font-size: 15px; line-height: 1.6; color: #555555;">%s</p>`, text)
}

func smallText(text string) string {
	return fmt.Sprintf(`<p style="margin: 16px 0 0; font-size: 12px; color: #999999; text-align: center;">%s</p>`, text)
}
