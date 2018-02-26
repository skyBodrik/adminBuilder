/**
 * Classe Url
 * @author Patrick Poulain
 * @see http://petitchevalroux.net
 * @licence GPL
 **/
function Url() {
}
Url.prototype.params = new Array();
Url.prototype.getQuery = function(url) {
    var str = url;
    var strpos = str.indexOf('?');
    /** Si on ne trouve pas de queryString on retourne une chaine vide */
    if (strpos == -1) {
        return '';
    }
    str = str.substr(strpos + 1, str.length);
    /**
     * Maintenant on verifie si on a une anchor ou pas (#) et si c'est le cas on
     * arrete la querystring avant
     */
    strpos = str.indexOf('#');
    if(strpos == -1)
    {
        return str;
    }
    return str.substr(0,strpos);
}
Url.prototype.getPath = function(url) {
    var strpos = url.indexOf('?');
    /** Si on ne trouve pas de queryString on retourne une chaine vide */
    if (strpos == -1) {
        return url;
    }
    return url.substr(0, strpos);
}

Url.prototype.buildParamFromString =  function(param)
{
    var p = decodeURIComponent(param);
    var strpos = p.indexOf('=');
    /*Si on trouve pas d'égale on met le paramètre a ''*/
    if(strpos == -1 )
    {
        this.params[p] = '';
        this.params.length++;
        return true;
    }
    var name = p.substr(0,strpos);
    var value = p.substr(strpos+1,p.length);
    var openBracket = name.indexOf('[');
    var closeBracket = name.indexOf(']');
    /**On traite les paramètre qui ne sont pas sous forme de tableau*/
    if(openBracket == -1 || closeBracket == -1)
    {
        if(!(openBracket == -1 && closeBracket == -1))
        {
            name = name.replace(new RegExp('[\\[\\]]'),'_');
        }
        this.params[name] = value;
        return true;
    }
    var matches = name.match(new RegExp('\\[.*?\\]','g'));
    name = name.substr(0,openBracket);
    p = 'this.params';
    var key = name;
    for(i in matches)
    {
        p += '[\''+key+'\']';
        if(eval(p) == undefined || typeof(eval(p)) != 'object')
        {
            eval(p +'= new Array();');
        }
        key = matches[i].substr(1,matches[i].length-2);
        /*si la clé est null on met la longueur du tableau*/
        if(key == '')
        {
            key = eval(p).length;
        }
    }
    p += '[\''+key+'\']';
    eval(p +'= \''+value+'\';');
}
Url.prototype.parseQuery = function(queryString) {

    var str = queryString;
    str = str.replace(new RegExp('&'), '&');
    this.params = new Array();
    this.params.length = 0;
    str = str.split('&');
    var p = '';
    var startPos = -1;
    var endPos = -1;
    var arrayName = '';
    var arrayKey = ''
    for ( var i = 0; i < str.length; i++) {
        this.buildParamFromString(str[i]);
    }
    return this.params;
}
Url.prototype.buildStringFromParam = function(object,prefix)
{
    var p = '';
    var value ='';
    if(prefix != undefined)
    {
        p = prefix;
    }
    if(typeof(object) == 'object')
    {
        for(var name in object){
            value = object[name];
            name = p == '' ? name : '['+name+']';
            if(typeof(value) == 'object')
            {
                this.buildStringFromParam(value,p+name);
            }
            else
            {
                this.params[this.params.length] = p+name+'='+value;
            }
        }
    }
}
Url.prototype.buildQuery = function(params) {
    this.params = new Array();
    this.buildStringFromParam(params);
    return this.params.join('&');
}
var Url = new Url();