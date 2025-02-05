package handler

import (
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service fundraise.Usecase
}

func New(service fundraise.Usecase) fundraise.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetFundraises() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		fundraises := ctl.service.FindAll(page, size)

		if fundraises == nil {
			return ctx.JSON(404, helper.Response("There is No Fundraises!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": fundraises,
		}))
	}
}

func (ctl *controller) FundraiseDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)
		if fundraise == nil {
			return ctx.JSON(404, helper.Response("Fundraise Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": fundraise,
		}))
	}
}

func (ctl *controller) CreateFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputFundraise{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		fundraise := ctl.service.Create(input)

		if fundraise == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": fundraise,
		}))
	}
}

func (ctl *controller) UpdateFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputFundraise{}

		fundraiseID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("Fundraise Not Found!"))
		}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, fundraiseID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Fundraise Success Updated!"))
	}
}

func (ctl *controller) DeleteFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("Fundraise Not Found!"))
		}

		delete := ctl.service.Remove(fundraiseID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Fundraise Success Deleted!", nil))
	}
}
