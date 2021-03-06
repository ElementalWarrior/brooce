package web

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"brooce/config"
	myredis "brooce/redis"
	"brooce/task"

	redis "gopkg.in/redis.v6"
)

type PagedHits struct {
	Hits       []*task.Task
	Start      int
	End        int
	Pages      int
	PageSize   int
	PageWanted int
}

func searchHandler(req *http.Request, rep *httpReply) (err error) {
	query, queue, listType, page := searchQueryParams(req.URL.RawQuery)

	key := fmt.Sprintf("%s:queue:%s:%s", redisHeader, queue, listType)

	hits := searchQueueForCommand(key, query)
	pagedHits := newPagedHits(hits, 10, page)

	output := &joblistOutputType{
		QueueName: queue,
		ListType:  listType,
		Query:     query,
		Page:      int64(page),
		Pages:     int64(pagedHits.Pages),
		Jobs:      pagedHits.Hits,
		Start:     int64(pagedHits.Start),
		End:       int64(pagedHits.End),
		Length:    int64(len(hits)),

		URL: req.URL,
	}

	err = templates.ExecuteTemplate(rep, "joblist", output)
	return
}

func newPagedHits(hits []*task.Task, pageSize int, pageWanted int) *PagedHits {
	if pageWanted < 1 {
		pageWanted = 1
	}

	start := 1
	end := pageSize

	totalHits := len(hits)
	totalPages := int(math.Ceil(float64(totalHits) / float64(pageSize)))

	maxStart := (pageWanted - 1) * pageSize
	maxEnd := (pageWanted * pageSize) - 1

	if maxStart > totalHits {
		start = totalHits
	} else {
		start = maxStart
	}

	if (maxEnd + 1) > totalHits {
		end = totalHits
	} else {
		end = maxEnd + 1
	}

	// log.Printf("page %d: start: %d end: %d total pages: %d", pageWanted, start, end, totalPages)

	return &PagedHits{Hits: hits[start:end], Start: start + 1, End: end, PageWanted: pageWanted, Pages: totalPages}
}

func searchQueryParams(rq string) (query string, queue string, listType string, page int) {
	params, err := url.ParseQuery(rq)
	if err != nil {
		log.Printf("Malformed URL query: %s err: %s", rq, err)
		return "", "", "done", 1
	}

	query = params.Get("q")
	queue = params.Get("queue")
	listType = params.Get("listType")
	if listType == "" {
		listType = "done"
	}

	page = 1
	if pg, err := strconv.Atoi(params.Get("page")); err == nil && pg > 1 {
		page = pg
	}

	return query, queue, listType, page
}

func searchQueueForCommand(queueKey string, query string) []*task.Task {
	// log.Printf("Searching %s for %s", queueKey, query)
	r := myredis.Get()

	found := []*task.Task{}
	vals := r.LRange(queueKey, 0, -1).Val()

	for _, v := range vals {
		// log.Printf("%s: %+v", queueKey, v)
		t, err := task.NewFromJson(v, config.JobOptions{})

		if err != nil {
			log.Printf("Couldn't construct task.Task from %+v", v)
			continue
		}

		if strings.Contains(t.Command, query) {
			found = append(found, t)
		}
	}

	// log.Printf("Search of %s for %s got %d hits", queueKey, query, len(found))

	found = addLogsToSearchHits(found)
	return found
}

func addLogsToSearchHits(hits []*task.Task) []*task.Task {
	hasLog := make([]*redis.IntCmd, len(hits))
	_, err := redisClient.Pipelined(func(pipe redis.Pipeliner) error {
		for i, job := range hits {
			k := fmt.Sprintf("%s:jobs:%s:log", redisHeader, job.Id)
			hasLog[i] = pipe.Exists(k)
		}
		return nil
	})
	if err != nil {
		return []*task.Task{}
	}

	for i, result := range hasLog {
		hits[i].HasLog = result.Val() > 0
	}

	return hits
}
