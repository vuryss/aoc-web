package service

import (
	"../core"
	"errors"
	"strconv"
	"time"
)

var allowedYears = map[int]bool{2015: true, 2016: true, 2017: true}

type YearIndexService struct {
	*core.Service
}

func (c *YearIndexService) Index() core.Response {
	// Validate year input
	year, err := c.validateYear()

	if err != nil {
		c.View.PassValue("layout_error", err.Error())
		return c.RedirectWithViewData("Index", "Index")
	}

	date, _ := time.Parse("2006-01-02", strconv.Itoa(year) + "-12-01")

	// How many dates to skip from the calendar
	weekday := date.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	c.View.PassValue("skipDays", int(weekday) - 1)

	// Pass year
	c.View.PassValue("year", year)

	c.View.SetView("year.html")

	return &core.ViewResponse{View: c.View}
}

func (c *YearIndexService) validateYear() (int, error) {
	yearString, isSet := c.Service.Parameters["year"]

	if !isSet {
		return 0, errors.New("missing year")
	}

	year, err := strconv.Atoi(yearString)

	if err != nil {
		return 0, errors.New("incorrect year")
	}

	if _, exists := allowedYears[year]; !exists {
		return 0, errors.New("year not supported")
	}

	return year, nil
}
