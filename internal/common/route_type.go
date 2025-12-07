package common

type RouteType string

const (
    RouteStatic       RouteType = "Static"
    RouteReverseProxy RouteType = "ReverseProxy"
    RouteRedirect     RouteType = "Redirect"
)


func (r RouteType) IsValid() bool {
    switch r {
    case RouteStatic, RouteReverseProxy, RouteRedirect:
        return true
    }
    return false
}


func (r RouteType) String() string {
    return string(r)
}