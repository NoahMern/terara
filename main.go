package main

import (
	"github.com/noahmern/terara/pkg/lexer"
)

func main() {
	l := lexer.NewLexer(`
	param($from_id,$to_id,$amount);
	use(ice);
	let balance = colletion::transfers.filter(id = $from_id).select('amount').sum();
	if(balance > $amount).
	then(collection::transfers.insert(
		document::new($from_id,$to_id,$amount).union(
		{'id': uuid(),
			'timestamp': now()})
	));
	`)
	for {
		token, err := l.NextToken()
		if err != nil {
			panic(err)
		}
		if token.Type == lexer.TokenEOF {
			break
		}
		println(token.String() + ": " + token.Value)
	}
}
