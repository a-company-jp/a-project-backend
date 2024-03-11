package handler

import (
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/user"
	"a-project-backend/svc/pkg/domain/query"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
)

type User struct {
	userQ query.User
}

func NewUser(userQ query.User) User {
	return User{
		userQ: userQ,
	}
}

func (u User) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := user.ID(c.Param("user_id"))
		_, err := u.userQ.GetUserByID(userID)
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			} else {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
		}
		// TODO: create response using protobuf struct
		resp := pb_out.UserInfoResponse{
			UserData: &pb_out.UserData{
				UserId:        "96366846-918a-45d9-a7e3-a382c0687ef4",
				Username:      "Shion1305",
				Firstname:     "詩恩",
				Lastname:      "市川",
				FirstnameKana: "Shion",
				LastnameKana:  "Ichikawa",
				StatusMessage: "Hello, World!",
				Tag: []*pb_out.Tag{
					{
						TagId:   "aa466caf-3596-4b67-a83b-d009a6dd2e91",
						TagName: "Security",
					},
					{
						TagId:   "2343486c-a54b-4629-9e8f-a25d6405faed",
						TagName: "DevOps",
					},
					{
						TagId:   "aa466caf-3596-4b67-a83b-d009a6dd2e91",
						TagName: "SRE",
					},
				},
				IconImageHash: "62063285-95c4-4ad2-85d5-dfe609706c46",
			},
			Milestones: []*pb_out.Milestone{
				{
					UserId:     "96366846-918a-45d9-a7e3-a382c0687ef4",
					EventId:    "86dd565a-57f1-4c49-b21c-ff3725db6109",
					Title:      "ペネトレーションテストができる専門家",
					Content:    "現代のデジタル化された社会において、セキュリティは極めて重要な要素です。組織、企業、さらには個人のデータも、常にサイバー攻撃の脅威に晒されています。このような環境で、ペネトレーションテスターとしての能力は、最前線でデジタルセキュリティを守る上で不可欠です。\n\nペネトレーションテストができるようになるためには、まず、ネットワークやシステムの仕組みを深く理解することから始まります。この知識を土台として、サイバー攻撃者が利用する可能性のある弱点を特定し、検証する技術を習得します。ペネトレーションテストでは、実際の攻撃を模倣して、セキュリティの脆弱性を発見し、報告することが求められます。このプロセスを通じて、組織はその脆弱性を修正し、セキュリティを強化することができます。\n\n技術的なスキルだけでなく、高い倫理観もこの職業には不可欠です。ペネトレーションテスターは、自分の行動がクライアントのセキュリティと信頼に直接影響を与えることを常に意識しておく必要があります。また、複雑なセキュリティ問題を解決するための創造的な思考も求められます。\n\n将来のペネトレーションテスターには、次のようなステップが推奨されます：\n\nコンピューターサイエンス、情報セキュリティ、または関連分野での教育を受ける。\nセキュリティ関連の資格（例：CEH、OSCP）を取得する。\n実践的な経験を積むために、セキュリティラボやハッキングコンテストに参加する。\n業界の動向に常に敏感であり続けるために、最新のセキュリティ技術や脅威について学び続ける。\nペネトレーションテスターとしてのキャリアは、技術的な挑戦と倫理的な責任を伴いますが、デジタル世界をより安全な場所にするための非常に重要な役割を果たします。この道を選ぶことは、社会に対する大きな貢献となり得ます。\n\n",
					ImageHash:  "7d4e7cdf-eaf5-476e-8800-d42b103a3678",
					BeginDate:  "2024-01-01",
					FinishDate: "2024-05-31",
				},
			},
		}
		respData, err := proto.Marshal(&resp)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/octet-stream", respData)
	}
}
