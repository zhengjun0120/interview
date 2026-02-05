package error_msg

import "errors"

var GET_USER_ID_ERROR = errors.New("获取用户ID失败")

var USER_ID_NOT_EXIST = errors.New("用户ID不存在")

var USER_ID_NOT_NULL = errors.New("用户ID不能为空")

var USER_NOT_EXIST = errors.New("用户不存在")

var USER_TYPE_ERROR = errors.New("用户类型错误")

var PASSWORD_NOT_NULL = errors.New("密码不能为空")

var USERNAME_NOT_NULL = errors.New("用户名不能为空")

var PASSWORD_ERROR = errors.New("密码错误或用户不存在")

var EMAIL_EXIST = errors.New("邮箱已注册")

var EMAIL_NOT_EXIST = errors.New("邮箱未注册")

var (
	RESUME_NAME_NOT_NULL            = errors.New("简历名称不能为空")
	RESUME_ID_NOT_NULL              = errors.New("简历ID不能为空")
	RESUME_NAME_TOO_LONG            = errors.New("简历名称太长")
	TALENT_ID_NOT_NULL              = errors.New("人才ID不能为空")
	TALENT_NAME_NOT_NULL            = errors.New("人才名称不能为空")
	TALENT_NAME_TOO_LONG            = errors.New("人才名称太长")
	TALENT_TARGET_POSITION_NOT_NULL = errors.New("职位不能为空")
	TALENT_MATCH_SCORE_INVALID      = errors.New("匹配分数无效")
)
