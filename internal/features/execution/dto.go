package execution

type ExecuteCommandRequest struct {
	Target  Target `json:"target" binding:"required"`
	Command string `json:"command" binding:"required"`
}

type Target struct {
	IP     string `json:"ip" binding:"required,ip"`
	Driver string `json:"driver" binding:"required"`
	Auth   Auth   `json:"auth" binding:"required"`
}

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Port     int    `json:"port"`
}

type ExecuteCommandResponse struct {
	Status string `json:"status"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type GetStatsRequest struct {
	Target Target `json:"target" binding:"required"`
}

type GetStatsResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}
