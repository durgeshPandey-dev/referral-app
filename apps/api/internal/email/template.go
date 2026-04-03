package email

import "fmt"

func BuildTemplate(name string, company string) string {
	if name == "" {
		name = "Hiring Team"
	}
	if company == "" {
		company = "your organization"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height:1.6; color:#333;">

<p>Hi %s,</p>

<p>
I hope you're doing well. I came across opportunities at <b>%s</b> and wanted to reach out.
</p>

<p>
I have <b>3.6 years of experience</b> working at <b>NPCI</b> and <b>EdgeVerve Systems</b>, where I’ve built 
<b>scalable backend systems</b> using <b>Golang</b> and <b>Node.js</b>.
</p>

<p>
My experience includes:
<ul>
<li>Designing and developing high-performance APIs</li>
<li>Building microservices for financial systems</li>
<li>Optimizing databases for high-throughput applications</li>
</ul>
</p>

<p>
I’m currently exploring new opportunities and would love to contribute to your team.  
I am available to join on <b>short notice</b>.
</p>

<p>
I’ve attached my resume for your reference.  
Would love to connect if there’s a relevant opportunity.
</p>

<p>
Best regards,<br>
<b>Durgesh Pande</b><br>
📧 durgeshpande20@gmail.com<br>
📞 +91 7348008896
</p>

</body>
</html>
`, name, company)
}
