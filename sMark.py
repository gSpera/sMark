from pygments.lexer import RegexLexer
from pygments.token import *

class sMarkLexer(RegexLexer):
    name = 'sMark'
    aliases = ['sm', 'smark']
    filenames = ['*.sm']
    mimetypes = ["text/x-sMark"]

    tokens = {
        'root': [
            (r' .*\n', Text),
            (r'^\s*.+\n\s*=+\n', Generic.Heading), #Title
            (r'^\s*.+\n\s*-+\n', Generic.Subheading), #Subtitle
            (r'^\|.*\n', Literal), #Quote
            (r'^-+ .*\n', Literal), #List
            (r'\[.\]', Generic.Strong), #Checkbox
            (r'\[.{2,}\]\n', Generic.Output, 'code'), #Code Block
            (r'\n-{2,}\n', Generic.Emph), #Divider
            (r'".+"@".+"', String.Single), #Link and Image
            (r'\*.*\*', Generic.Strong), #Bold
            (r'/.*/', Generic.Emph), #Italic
            (r'_.*_', Generic.Emph), #Underline
            (r'\-.*\-', Generic.Emph), #Strikethrought
            (r'\++\n(.*\n)+\++\n', Generic.Output), #Header
            (r'.*\n', Text),
        ],
        'code': [
            (r'.*\n', Generic.Output),
            (r'\[end\]\n', Generic.Output, '#pop'),
        ],
    }