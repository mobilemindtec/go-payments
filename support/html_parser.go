package support

import (
    "golang.org/x/net/html"
    "strings"
    _ "errors"
    "bytes"
    _ "fmt"
    "io"
)

func HtmlParse(htmlContent string) (*html.Node, error){
    return html.Parse(strings.NewReader(htmlContent))
}

func HtmlParseFindByName(doc *html.Node, elementName string) *html.Node {
    var body *html.Node
    var crawler func(*html.Node)
    crawler = func(node *html.Node) {
        if node.Type == html.ElementNode {

            for _, a := range node.Attr {

                if a.Key == "name" && a.Val == elementName {
                    //fmt.Println("html key = %v, val = %v", a.Key, a.Val)
                    body = node
                    return
                }
            }

        }
        for child := node.FirstChild; child != nil; child = child.NextSibling {
            crawler(child)
        }
    }
    crawler(doc)
    return body
}

func HtmlParseFindById(doc *html.Node, elementId string) *html.Node {
    var body *html.Node
    var crawler func(*html.Node)
    crawler = func(node *html.Node) {
        if node.Type == html.ElementNode {

            for _, a := range node.Attr {

                if a.Key == "id" && a.Val == elementId {
                    //fmt.Println("html key = %v, val = %v", a.Key, a.Val)
                    body = node
                    return
                }
            }

        }
        for child := node.FirstChild; child != nil; child = child.NextSibling {
            crawler(child)
        }
    }
    crawler(doc)
    return body
}

func HtmlParseFindByType(doc *html.Node, elementType string) *html.Node {
    var body *html.Node
    var crawler func(*html.Node)
    var first = true
    crawler = func(node *html.Node) {
       if !first {          
            if node.Type == html.ElementNode && node.Data == elementType {
                body = node
                return
            }
        }
       
        first = false
        
        for child := node.FirstChild; child != nil; child = child.NextSibling {
            crawler(child)
        }
    }
    
    crawler(doc)
    return body
}

func NodeToString(n *html.Node) string {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    html.Render(w, n)
    return buf.String()
}

func GetNodeAttrValue(n *html.Node, attrName string) string {

       for _, a := range n.Attr {
            if a.Key == attrName {
        return a.Val
            }
        }
    return ""

}