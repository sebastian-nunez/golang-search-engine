package core

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

// createIndexTokens processes the input text through a series of text normalization steps and returns the resulting tokens
// which can be used for tasks like indexing, searching, or text classification. The analysis pipeline consists of the
// following stages:
//
//  1. Tokenization: splits the text into individual words or tokens.
//  2. Lowercasing: converts all tokens to lowercase for case-insensitive
//     matching.
//  3. Stop word removal: removes common words (e.g., "the," "a," "is") that
//     usually don't carry significant meaning.
//  4. Stemming: reduces words to their root form (e.g., "running" becomes
//     "run").
func createIndexTokens(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)

	return tokens
}

// tokenize returns a slice of tokens for the given text.
// A token is defined as a contiguous sequence of letters or numbers.
// Any character that is not a letter or a number is considered a separator and is ignored.
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Defines what the separator is
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// lowercaseFilter returns a slice of tokens normalized to lower case.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, t := range tokens {
		r[i] = strings.ToLower(t)
	}

	return r
}

// stopwordFilter returns a slice of tokens with stop words removed.
func stopwordFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for _, t := range tokens {
		_, ok := stopwords[t]
		if !ok {
			r = append(r, t)
		}
	}

	return r
}

// stemmerFilter applies stemming to a slice of tokens using the Snowball English stemmer.
// Stemming is a text normalization technique that reduces words to their root form (or "stem").
// For example, "running," "runs," and "run" would all be stemmed to "run." This process helps
// improve search recall by matching different forms of the same word.
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, t := range tokens {
		r[i] = snowballeng.Stem(t, false)
	}

	return r
}

// stopwords is a set containing English words which should be
// removed from the text before indexing or analysis.
var stopwords = map[string]struct{}{
	"a": {}, "about": {}, "above": {}, "across": {}, "after": {},
	"afterwards": {}, "again": {}, "against": {}, "all": {}, "am": {},
	"among": {}, "amongst": {}, "amoungst": {}, "amount": {}, "an": {}, "and": {},
	"another": {}, "any": {}, "anyhow": {}, "anyone": {}, "anything": {}, "anyway": {},
	"anywhere": {}, "are": {}, "around": {}, "as": {}, "at": {}, "back": {}, "be": {},
	"became": {}, "because": {}, "become": {}, "becomes": {}, "becoming": {}, "been": {},
	"before": {}, "beforehand": {}, "behind": {}, "being": {}, "below": {}, "beside": {},
	"besides": {}, "between": {}, "beyond": {}, "bill": {}, "both": {},
	"bottom": {}, "but": {}, "by": {}, "call": {}, "can": {}, "cannot": {}, "cant": {},
	"co": {}, "con": {}, "could": {}, "couldn't": {}, "cry": {}, "de": {}, "describe": {},
	"detail": {}, "do": {}, "done": {}, "down": {}, "downwards": {}, "during": {}, "each": {},
	"eg": {}, "eight": {}, "either": {}, "eleven": {}, "else": {}, "elsewhere": {}, "empty": {},
	"enough": {}, "etc": {}, "even": {}, "ever": {}, "every": {}, "everyone": {},
	"everything": {}, "everywhere": {}, "except": {}, "few": {}, "fifteen": {},
	"fill": {}, "find": {}, "fire": {}, "first": {}, "five": {}, "for": {}, "former": {},
	"formerly": {}, "forty": {}, "found": {}, "four": {}, "from": {}, "front": {}, "full": {},
	"further": {}, "furthermore": {}, "get": {}, "give": {}, "go": {}, "had": {}, "has": {},
	"hasn't": {}, "have": {}, "he": {}, "hence": {}, "her": {}, "here": {}, "hereafter": {},
	"hereby": {}, "herein": {}, "hereupon": {}, "hers": {}, "herself": {}, "him": {},
	"himself": {}, "his": {}, "how": {}, "however": {}, "hundred": {}, "ie": {}, "if": {},
	"in": {}, "inc": {}, "indeed": {}, "interest": {}, "into": {}, "is": {}, "it": {}, "its": {},
	"itself": {}, "keep": {}, "last": {}, "latter": {}, "latterly": {}, "least": {}, "less": {},
	"lest": {}, "let": {}, "like": {}, "ltd": {}, "made": {}, "make": {}, "many": {}, "may": {},
	"me": {}, "meanwhile": {}, "might": {}, "mill": {}, "mine": {}, "more": {}, "moreover": {},
	"most": {}, "mostly": {}, "move": {}, "much": {}, "must": {}, "my": {}, "myself": {},
	"name": {}, "namely": {}, "neither": {}, "never": {}, "nevertheless": {}, "next": {},
	"nine": {}, "no": {}, "nobody": {}, "none": {}, "noone": {}, "nor": {}, "not": {},
	"nothing": {}, "now": {}, "nowhere": {}, "of": {}, "off": {}, "often": {}, "on": {},
	"once": {}, "one": {}, "only": {}, "onto": {}, "or": {}, "other": {}, "others": {},
	"otherwise": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {}, "outside": {}, "over": {},
	"own": {}, "part": {}, "per": {}, "perhaps": {}, "please": {}, "put": {}, "rather": {},
	"re": {}, "same": {}, "see": {}, "seem": {}, "seemed": {}, "seeming": {}, "seems": {},
	"several": {}, "she": {}, "should": {}, "show": {}, "side": {}, "since": {}, "sincere": {},
	"six": {}, "sixty": {}, "so": {}, "some": {}, "somehow": {}, "someone": {},
	"something": {}, "sometime": {}, "sometimes": {}, "somewhat": {}, "somewhere": {},
	"still": {}, "such": {}, "system": {}, "take": {}, "ten": {}, "than": {}, "that": {},
	"the": {}, "their": {}, "them": {}, "themselves": {}, "then": {}, "thence": {}, "there": {},
	"thereafter": {}, "thereby": {}, "therefore": {}, "therein": {}, "thereupon": {}, "these": {},
	"they": {}, "thickv": {}, "thin": {}, "third": {}, "this": {}, "those": {}, "though": {},
	"three": {}, "through": {}, "throughout": {}, "thru": {}, "thus": {}, "to": {},
	"together": {}, "too": {}, "top": {}, "toward": {}, "towards": {}, "twelve": {},
	"twenty": {}, "two": {}, "un": {}, "under": {}, "until": {}, "up": {}, "upon": {}, "us": {},
	"very": {}, "via": {}, "was": {}, "we": {}, "well": {}, "were": {}, "what": {},
	"whatever": {}, "when": {}, "whence": {}, "whenever": {}, "where": {}, "whereafter": {},
	"whereas": {}, "whereby": {}, "wherein": {}, "whereupon": {}, "wherever": {}, "whether": {},
	"which": {}, "while": {}, "whither": {}, "who": {}, "whoever": {}, "whole": {}, "whom": {},
	"whose": {}, "why": {}, "will": {}, "with": {}, "within": {}, "without": {}, "would": {},
	"yet": {}, "you": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
	"www": {}, "com": {}, "org": {}, "net": {}, "io": {}, "https": {}, "http": {},
	"html": {}, "php": {}, "asp": {},
}
