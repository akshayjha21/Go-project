package pagination

type Paginate struct {
   Limit int
   Page  int
}

func NewPaginate(limit int, page int) *Paginate {
   if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }
    return &Paginate{Limit: limit, Page: page}
}

func (p *Paginate) LimitOffset()(limit int, offset int) {
   offset = (p.Page - 1) * p.Limit

   return p.Limit,offset
}