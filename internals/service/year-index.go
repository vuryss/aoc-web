package service

import (
	"../core"
	"strconv"
)

var allowedYears = map[int]bool{2015: true, 2016: true, 2017: true}

type YearIndexService struct {
	*core.Service
}

func (c *YearIndexService) Index() core.Response {
	// Validate year input
	yearString, isSet := c.Service.Parameters["year"]

	if !isSet {
		c.View.PassValue("layout_error", "Year not found!")
		return c.RedirectWithViewData("Index", "Index")
	}

	year, err := strconv.Atoi(yearString)

	if err != nil {
		c.View.PassValue("layout_error", "Incorrect year value!")
		return c.RedirectWithViewData("Index", "Index")
	}

	if _, exists := allowedYears[year]; !exists {
		c.View.PassValue("layout_error", "Year not supported")
		return c.RedirectWithViewData("Index", "Index")
	}

	c.View.SetView("year.html")

	return &core.ViewResponse{View: c.View}
}
