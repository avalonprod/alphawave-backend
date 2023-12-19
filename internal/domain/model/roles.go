package model

// ROLES CONSTANT
const (
	ROLE_ADMIN       = "ADMIN"
	ROLE_OWNER       = "OWNER"
	ROLE_MANAGER     = "MANAGER"
	ROLE_EDITOR      = "EDITOR"
	ROLE_VIEWER      = "VIEWER"
	ROLE_AUDITOR     = "AUDITOR"
	ROLE_CONTRIBUTOR = "CONTRIBUTOR"
)

// PERMISSIONS CONSTANT
const (
	PERMISSION_UPLOAD_FILES           = "UPLOAD_FILES"
	PERMISSION_CREATE_PROJECTS        = "CREATE_PROJECTS"
	PERMISSION_DELETE_PROJECTS        = "DELETE_PROJECTS"
	PERMISSION_DELETE_FILES           = "DELITE_FILES"
	PERMISSION_SHARE_FILES_INTERNALLY = "SHARE_FILES_INTERNALLY"
	PERMISSION_SHARE_FILES_EXTERNALLY = "SHARE_FILES_EXTERNALLY"
	PERMISSION_ACCESS_AI_CHAT         = "ACCESS_AI_CHAT"
	PERMISSION_EDIT_PROJECT_PAGE      = "EDIT_PROJECT_PAGE"
	PERMISSION_ADD_NEW_TASKS          = "ADD_NEW_TASKS"
	PERMISSION_ACCESS_MESSAGER        = "ACCESS_MESSAGER"
	PERMISSION_INVITE_MEMBERS         = "INVITE_MEMBERS"
	PERMISSION_EDIT_JOB_TITLE         = "EDIT_JOB_TITLE"
	PERMISSION_EDIT_NAME              = "EDIT_NAME"
	PERMISSION_EDIT_EMAIL_ADDRESS     = "EDIT_EMAIL_ADDRESS"
	PERMISSION_INVITE_ADMINS          = "INVITE_ADMINS"
)

type TeamRoles struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	TeamID string `json:"teamID" bson:"teamID"`
	Roles  []Role `json:"roles" bson:"roles"`
}

type Role struct {
	Role        string      `json:"role" bson:"role"`
	Permissions Permissions `json:"permissions" bson:"permissions"`
}

type Permissions struct {
	UploadFiles          bool
	CreateProjects       bool
	DeleteProjects       bool
	DeleteFiles          bool
	ShareFilesInternally bool
	ShareFilesExternally bool
	AccessAiChat         bool
	EditProjectPage      bool
	AddNewTasks          bool
	AccessMessager       bool
	InviteMembers        bool
	EditJobTitle         bool
	EditName             bool
	EditEmailAddress     bool
	InviteAdmins         bool
}

type UpdatePermissions struct {
	Role        string
	Permissions Permissions
}

func CreateRoles(teamID string) TeamRoles {
	// defaultPermissions := Permissions{
	// 	UploadFiles: true,
	// }

	adminRole := Role{
		Role: ROLE_ADMIN,
	}
	ownerRole := Role{
		Role: ROLE_OWNER,
		Permissions: Permissions{
			UploadFiles:          true,
			CreateProjects:       true,
			DeleteProjects:       true,
			DeleteFiles:          true,
			ShareFilesInternally: true,
			ShareFilesExternally: true,
			AccessAiChat:         true,
			EditProjectPage:      true,
			AddNewTasks:          true,
			AccessMessager:       true,
			InviteMembers:        true,
			EditJobTitle:         true,
			EditName:             true,
			EditEmailAddress:     true,
			InviteAdmins:         true,
		},
	}
	managerRole := Role{
		Role: ROLE_MANAGER,
	}
	editorRole := Role{
		Role: ROLE_EDITOR,
	}
	viewerRole := Role{
		Role: ROLE_VIEWER,
	}
	auditorRole := Role{
		Role: ROLE_AUDITOR,
	}
	contributorRole := Role{
		Role: ROLE_CONTRIBUTOR,
	}

	var roles []Role

	roles = append(roles, adminRole, ownerRole, managerRole, editorRole, viewerRole, auditorRole, contributorRole)

	return TeamRoles{
		TeamID: teamID,
		Roles:  roles,
	}
}

func GetAvailableRoles() []Role {
	adminRole := Role{
		Role: ROLE_ADMIN,
	}
	ownerRole := Role{
		Role: ROLE_OWNER,
	}
	managerRole := Role{
		Role: ROLE_MANAGER,
	}
	editorRole := Role{
		Role: ROLE_EDITOR,
	}
	viewerRole := Role{
		Role: ROLE_VIEWER,
	}
	auditorRole := Role{
		Role: ROLE_AUDITOR,
	}
	contributorRole := Role{
		Role: ROLE_CONTRIBUTOR,
	}

	var roles []Role

	roles = append(roles, adminRole, ownerRole, managerRole, editorRole, viewerRole, auditorRole, contributorRole)
	return roles
}

func IsAvailableRole(role string) bool {
	availableRoles := GetAvailableRoles()

	for _, item := range availableRoles {
		if role == item.Role {
			return true
		}
	}
	return false
}
