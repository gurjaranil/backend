package utill

import "fmt"

func GenerateInvitationEmailHTML(recipientName, registrationLink, libraryEmail, libraryName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Library Registration Invitation</title>
  <style>
    /* Styles for the email template */
    body {
      font-family: Arial, sans-serif;
      line-height: 1.6;
      background-color: #f4f4f4;
      margin: 0;
      padding: 0;
    }

    .container {
      max-width: 600px;
      margin: 0 auto;
      padding: 20px;
      background-color: #ffffff;
      box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
      border-radius: 10px;
    }

    h1, h2 {
      color: #333333;
    }

    p {
      margin-bottom: 20px;
    }

    .btn {
      display: inline-block;
      padding: 10px 20px;
      background-color: #007bff;
      color: #ffffff;
      text-decoration: none;
      border-radius: 5px;
    }

    .footer {
      margin-top: 30px;
      text-align: center;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Welcome to Our Library!</h1>
    <p>Dear %s,</p>
    <p>We are excited to invite you to register for our library and explore a world of knowledge, imagination, and inspiration.</p>
    <p>Benefits of library membership:</p>
    <ul>
      <li>Access to a vast collection of books, magazines, and digital resources.</li>
      <li>Exclusive events, workshops, and book clubs.</li>
      <li>Personalized recommendations based on your interests.</li>
      <li>And much more!</li>
    </ul>
    <p>Click the button below to start your registration:</p>
    <a href="%s" class="btn">Register Now</a>
    <p>If you have any questions or need assistance, please feel free to contact us at %s.</p>
    <div class="footer">
      <p>Thank you for choosing %s.</p>
      <p>&copy; 2024 %s. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
`, recipientName, registrationLink, libraryEmail, libraryName, libraryName)
}
