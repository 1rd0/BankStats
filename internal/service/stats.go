package service

import (
	"bankstats/internal/domain"
	"bankstats/pkg/pkg/cbrclient"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Client       *cbrclient.Client
	RequestPause time.Duration
}

func (s *Service) CollectRange(ctx context.Context, start, end time.Time, codes map[string]bool) ([]domain.RatePoint, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end before start")
	}
	days := int(end.Sub(start).Hours()/24) + 1
	var all []domain.RatePoint

	for i := 0; i < days; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		d := start.AddDate(0, 0, i)
		ddmmyyyy := d.Format("02/01/2006")

		v, err := s.Client.FetchByDate(ctx, ddmmyyyy)
		if err != nil {
			return nil, fmt.Errorf("fetch %s: %w", ddmmyyyy, err)
		}
		respDate, err := time.Parse("02.01.2006", v.Date)
		if err != nil {
			respDate = d
		}

		for _, val := range v.Valute {
			code := strings.ToUpper(strings.TrimSpace(val.CharCode))
			if len(codes) > 0 && !codes[code] {
				continue
			}
			raw := strings.ReplaceAll(strings.TrimSpace(val.Value), ",", ".")
			num, err := strconv.ParseFloat(raw, 64)
			if err != nil || val.Nominal == 0 {
				continue
			}
			perUnit := num / float64(val.Nominal)
			all = append(all, domain.RatePoint{
				Date:     respDate,
				CharCode: code,
				Name:     val.Name,
				PerUnit:  perUnit,
			})
		}

		if s.RequestPause > 0 && i+1 < days {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(s.RequestPause):
			}
		}
	}
	return all, nil
}

func (s *Service) CalcStats(points []domain.RatePoint, days int) domain.Stats {
	var st domain.Stats
	st.Days = days
	if len(points) == 0 {
		return st
	}

	first := true
	var sum float64
	for _, p := range points {
		sum += p.PerUnit
		if first {
			st.Max = domain.Extremum{Value: p.PerUnit, CharCode: p.CharCode, Name: p.Name, Date: p.Date}
			st.Min = st.Max
			first = false
			continue
		}
		if p.PerUnit > st.Max.Value {
			st.Max = domain.Extremum{Value: p.PerUnit, CharCode: p.CharCode, Name: p.Name, Date: p.Date}
		}
		if p.PerUnit < st.Min.Value {
			st.Min = domain.Extremum{Value: p.PerUnit, CharCode: p.CharCode, Name: p.Name, Date: p.Date}
		}
	}

	st.Average = sum / float64(len(points))
	st.Count = len(points)
	return st
}

func (s *Service) CalcStatsUsingAllAverage(selected, all []domain.RatePoint, days int) domain.Stats {
	st := s.CalcStats(selected, days)
	st.Average = calcAverage(all)
	return st
}

func calcAverage(points []domain.RatePoint) float64 {
	if len(points) == 0 {
		return 0
	}
	var sum float64
	for _, p := range points {
		sum += p.PerUnit
	}
	return sum / float64(len(points))
}
