package templates

type HelmBasic struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type HelmImage struct {
	Repository string `json:"repository"`
	PullPolicy string `json:"pullPolicy"`
	Tag        string `json:"tag"`
}

type HelmAccount struct {
	Name string `json:"name"`
}

type HelmContainer struct {
	TargetPort int `json:"targetPort"`
	Port       int `json:"port"`
}

type HelmHost struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type HelmTls struct {
	Hosts []string `json:"hosts"`
}

type HelmIngress struct {
	Hosts []*HelmHost `json:"hosts"`
	Tls   []*HelmTls  `json:"tls"`
}

type HelmMount struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	SubPath   string `json:"subPath"`
}

type HelmConfigData struct {
	Appsettings   string `json:"appsettings.json"`
	Log4netConfig string `json:"log4net.config"`
}

//goland:noinspection ALL
type HelmConfig struct {
	ConfigName    string          `json:"configName"`
	ConfigMapName string          `json:"configMapName"`
	Data          *HelmConfigData `json:"data"`
}

//goland:noinspection ALL
type ApiValuesDto struct {
	Body           *HelmBasic             `json:"blockraiders"`
	Count          int                    `json:"replicaCount"`
	Image          *HelmImage             `json:"image"`
	ServiceAccount *HelmAccount           `json:"serviceAccount"`
	ServiceType    string                 `json:"serviceType"`
	Container      *HelmContainer         `json:"container"`
	Ingress        *HelmIngress           `json:"ingress"`
	VolumeMounts   []*HelmMount           `json:"volumeMounts"`
	Config         *HelmConfig            `json:"config"`
	NodeAffinity   map[string]interface{} `json:"nodeAffinity"`
	Resources      map[string]interface{} `json:"resources"`
}

type GameHostValuesDto struct {
	Body           *HelmBasic             `json:"blockraiders"`
	Image          *HelmImage             `json:"image"`
	ServiceAccount *HelmAccount           `json:"serviceAccount"`
	Container      *HelmContainer         `json:"container"`
	VolumeMounts   []*HelmMount           `json:"volumeMounts"`
	Config         *HelmConfig            `json:"config"`
	NodeAffinity   map[string]interface{} `json:"nodeAffinity"`
	Resources      map[string]interface{} `json:"resources"`
}
