package kubectl

const TemplateValue = ":hostname"

func GetTemplate() string {
	return `apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: '` + TemplateValue + `'
  namespace: mpmain-frontend
  annotations:
    kubernetes.io/ingress.class: traefik
    traefik.ingress.kubernetes.io/frontend-entry-points: http,https
    traefik.ingress.kubernetes.io/redirect-permanent: "true"
    traefik.ingress.kubernetes.io/preserve-host: "true"
    traefik.ingress.kubernetes.io/redirect-entry-point: https
    ingress.kubernetes.io/ssl-redirect: "true"
spec:
  rules:
  - host: '` + TemplateValue + `'
    http:
      paths:
      - backend:
          serviceName: mpmain-frontend
          servicePort: 80
  - host: 'www.` + TemplateValue + `'
    http:
      paths:
      - backend:
          serviceName: mpmain-frontend
          servicePort: 80
  tls:
    - secretName: polecato-cert`
}
