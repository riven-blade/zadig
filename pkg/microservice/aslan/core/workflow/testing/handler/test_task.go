/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/koderover/zadig/v2/pkg/microservice/aslan/config"
	commonmodels "github.com/koderover/zadig/v2/pkg/microservice/aslan/core/common/repository/models"
	"github.com/koderover/zadig/v2/pkg/microservice/aslan/core/common/service/notify"
	workflowservice "github.com/koderover/zadig/v2/pkg/microservice/aslan/core/workflow/service/workflow"
	"github.com/koderover/zadig/v2/pkg/microservice/aslan/core/workflow/testing/service"
	"github.com/koderover/zadig/v2/pkg/setting"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
	e "github.com/koderover/zadig/v2/pkg/tool/errors"
	"github.com/koderover/zadig/v2/pkg/tool/log"
)

func CreateTestTask(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	args := new(commonmodels.TestTaskArgs)
	data, err := c.GetRawData()
	if err != nil {
		log.Errorf("CreateTestTask c.GetRawData() err : %v", err)
	}
	if err = json.Unmarshal(data, args); err != nil {
		log.Errorf("CreateTestTask json.Unmarshal err : %v", err)
	}
	projectKey := args.ProductName
	internalhandler.InsertOperationLog(c, ctx.UserName, projectKey, "新增", "测试-task", fmt.Sprintf("%s-%s", args.TestName, "job"), string(data), ctx.Logger)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.Execute {
			ctx.UnAuthorized = true
			return
		}
	}

	if err := c.ShouldBindWith(&args, binding.JSON); err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(err.Error())
		return
	}

	if args.TestTaskCreator != setting.CronTaskCreator && args.TestTaskCreator != setting.WebhookTaskCreator {
		args.TestTaskCreator = ctx.UserName
	}

	ctx.Resp, ctx.Err = service.CreateTestTaskV2(args, ctx.UserName, ctx.Account, ctx.UserID, ctx.Logger)
	if ctx.Err != nil {
		notify.SendFailedTaskMessage(ctx.UserName, args.ProductName, args.TestName, ctx.RequestID, config.TestType, ctx.Err, ctx.Logger)
	}
}

func ListTestTask(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	pageSizeStr := c.Query("pageSize")
	pageNumStr := c.Query("pageNum")

	var pageSize, pageNum int

	if pageSizeStr == "" {
		pageSize = 50
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("pageSize args err :%s", err))
			return
		}
	}

	if pageNumStr == "" {
		pageNum = 1
	} else {
		pageNum, err = strconv.Atoi(pageNumStr)
		if err != nil {
			ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("page args err :%s", err))
			return
		}
	}

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.View {
			ctx.UnAuthorized = true
			return
		}
	}

	ctx.Resp, ctx.Err = service.ListTestTask(testName, projectKey, pageNum, pageSize, ctx.Logger)
}

func CancelTestTaskV3(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	taskIDStr := c.Query("taskID")

	internalhandler.InsertOperationLog(c, ctx.UserName, projectKey, "取消", "测试任务", c.Param("name"), "", ctx.Logger)

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.Execute {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("taskID args err :%s", err))
		return
	}

	workflowName := fmt.Sprintf(setting.TestWorkflowNamingConvention, testName)

	ctx.Err = workflowservice.CancelWorkflowTaskV4(ctx.UserName, workflowName, taskID, ctx.Logger)
}

func GetTestTaskInfo(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	taskIDStr := c.Query("taskID")

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.View {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("taskID args err :%s", err))
		return
	}

	ctx.Resp, ctx.Err = service.GetTestTaskDetail(projectKey, testName, taskID, ctx.Logger)
}

func GetTestTaskReportInfo(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	taskIDStr := c.Query("taskID")

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.View {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("taskID args err :%s", err))
		return
	}

	ctx.Resp, ctx.Err = service.GetTestTaskReportDetail(projectKey, testName, taskID, ctx.Logger)
}

func RestartTestTaskV2(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	taskIDStr := c.Query("taskID")

	internalhandler.InsertOperationLog(c, ctx.UserName, projectKey, "重启", "测试任务", c.Param("name"), "", ctx.Logger)

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.Execute {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("taskID args err :%s", err))
		return
	}

	workflowName := fmt.Sprintf(setting.TestWorkflowNamingConvention, testName)

	ctx.Err = workflowservice.RetryWorkflowTaskV4(workflowName, taskID, ctx.Logger)
}

func GetTestingTaskArtifact(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Query("testName")
	jobName := c.Query("jobName")
	taskIDStr := c.Query("taskID")

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.View {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc("invalid task id")
		return
	}

	workflowName := fmt.Sprintf(setting.TestWorkflowNamingConvention, testName)

	resp, err := workflowservice.GetWorkflowV4ArtifactFileContent(workflowName, jobName, taskID, ctx.Logger)
	if err != nil {
		ctx.Err = err
		return
	}
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="artifact.tar.gz"`)

	c.Data(200, "application/octet-stream", resp)
}

func GetTestingTaskSSE(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	projectKey := c.Query("projectName")
	testName := c.Param("testName")
	taskIDStr := c.Param("taskID")

	// authorization check
	if !ctx.Resources.IsSystemAdmin {
		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
			ctx.UnAuthorized = true
			return
		}

		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
			!ctx.Resources.ProjectAuthInfo[projectKey].Test.View {
			ctx.UnAuthorized = true
			return
		}
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.Err = e.ErrInvalidParam.AddDesc(fmt.Sprintf("taskID args err :%s", err))
		return
	}

	internalhandler.Stream(c, func(ctx1 context.Context, msgChan chan interface{}) {
		startTime := time.Now()
		err := wait.PollImmediateUntil(time.Second, func() (bool, error) {
			res, err := service.GetTestTaskDetail(projectKey, testName, taskID, ctx.Logger)
			if err != nil {
				ctx.Logger.Errorf("[%s] GetWorkflowTaskV3SSE error: %s", ctx.UserName, err)
				return false, err
			}

			msgChan <- res

			if time.Since(startTime).Minutes() == float64(60) {
				ctx.Logger.Warnf("[%s] Query GetWorkflowTaskV3SSE API over 60 minutes", ctx.UserName)
			}

			return false, nil
		}, ctx1.Done())

		if err != nil && err != wait.ErrWaitTimeout {
			ctx.Logger.Error(err)
		}
	}, ctx.Logger)
}

// TODO: Deprecated code, remove after v2.2.0
//func RestartTestTask(c *gin.Context) {
//	ctx, err := internalhandler.NewContextWithAuthorization(c)
//	defer func() { internalhandler.JSONResponse(c, ctx) }()
//
//	if err != nil {
//		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
//		ctx.UnAuthorized = true
//		return
//	}
//
//	projectKey := c.Param("productName")
//
//	internalhandler.InsertOperationLog(c, ctx.UserName, projectKey, "重启", "测试任务", c.Param("name"), "", ctx.Logger)
//
//	// authorization check
//	if !ctx.Resources.IsSystemAdmin {
//		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
//			ctx.UnAuthorized = true
//			return
//		}
//
//		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
//			!ctx.Resources.ProjectAuthInfo[projectKey].Test.Execute {
//			ctx.UnAuthorized = true
//			return
//		}
//	}
//
//	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
//	if err != nil {
//		ctx.Err = e.ErrInvalidParam.AddDesc("invalid task id")
//		return
//	}
//
//	ctx.Err = workflowservice.RestartPipelineTaskV2(ctx.UserName, taskID, c.Param("name"), config.TestType, ctx.Logger)
//}
//
//func CancelTestTaskV2(c *gin.Context) {
//	ctx, err := internalhandler.NewContextWithAuthorization(c)
//	defer func() { internalhandler.JSONResponse(c, ctx) }()
//
//	if err != nil {
//
//		ctx.Err = fmt.Errorf("authorization Info Generation failed: err %s", err)
//		ctx.UnAuthorized = true
//		return
//	}
//
//	projectKey := c.Param("productName")
//	internalhandler.InsertOperationLog(c, ctx.UserName, projectKey, "取消", "测试任务", c.Param("name"), "", ctx.Logger)
//
//	// authorization check
//	if !ctx.Resources.IsSystemAdmin {
//		if _, ok := ctx.Resources.ProjectAuthInfo[projectKey]; !ok {
//			ctx.UnAuthorized = true
//			return
//		}
//
//		if !ctx.Resources.ProjectAuthInfo[projectKey].IsProjectAdmin &&
//			!ctx.Resources.ProjectAuthInfo[projectKey].Test.Execute {
//			ctx.UnAuthorized = true
//			return
//		}
//	}
//
//	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
//	if err != nil {
//		ctx.Err = e.ErrInvalidParam.AddDesc("invalid task id")
//		return
//	}
//	ctx.Err = commonservice.CancelTaskV2(ctx.UserName, c.Param("name"), taskID, config.TestType, ctx.RequestID, ctx.Logger)
//}
