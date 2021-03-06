package main

import (
	"time"
)

type JaxbPackageInfo struct {
	BaseNameSpace       string
	AdditionalNameSpace []*FQNAbbr
	PackageName         string
}

type JaxbMainClassInfo struct {
	PackageName       string
	BaseXMLClassName  string
	SourceXMLFilename string
	Date              time.Time
}

type JaxbClassInfo struct {
	Name                   string
	Root                   bool
	PackageName, ClassName string
	Attributes             []*JaxbAttribute
	Fields                 []*JaxbField
	HasValue               bool
	ValueType              string
	Date                   time.Time
}

type JaxbAttribute struct {
	Name      string
	NameUpper string
	NameLower string
	NameSpace string
}
type JaxbField struct {
	TypeName  string
	Name      string
	NameUpper string
	NameLower string
	NameSpace string
	Repeats   bool
}

func (jb *JaxbClassInfo) init() {
	jb.Attributes = make([]*JaxbAttribute, 0)
	jb.Fields = make([]*JaxbField, 0)
}

const jaxbClassTemplate = `
// Generated by chidley https://github.com/gnewton/chidley
// Date: {{.Date}}
// 
package {{.PackageName}}.xml;

import java.util.ArrayList;
import javax.xml.bind.annotation.*;
import com.google.gson.annotations.SerializedName;

@XmlAccessorType(XmlAccessType.FIELD)
@XmlRootElement(name="{{.Name}}")
public class {{.ClassName}} {
{{if .Attributes}}
    // Attributes{{end}}
{{range .Attributes}}
{{if .NameSpace}}    
@XmlAttribute(namespace = "{{.NameSpace}}"){{else}}    @XmlAttribute(name="{{.Name}}"){{end}}
    @SerializedName("{{.Name}}")
    public String {{.NameLower}};{{end}}
{{if .Fields}}
    // Fields{{end}}{{range .Fields}}    
    @XmlElement(name="{{.Name}}")
    @SerializedName("{{.Name}}")
    {{if .Repeats}}public ArrayList<{{.TypeName}}> {{.NameLower}}{{else}}public {{.TypeName}} {{.NameLower}}{{end}};
{{end}}
{{if .HasValue}}
    // Value
    @XmlValue
    public {{.ValueType}} tagValue;{{end}}
}
`

const jaxbMainTemplate = `
// Generated by chidley https://github.com/gnewton/chidley
// Date: {{.Date}}
//

package {{.PackageName}};
 
import java.io.File;
import javax.xml.bind.JAXBContext;
import javax.xml.bind.JAXBException;
import javax.xml.bind.Unmarshaller;
import {{.PackageName}}.xml.{{.BaseXMLClassName}};
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import java.net.URL;
import java.net.URLClassLoader;


 
public class Main {
	public static void main(String[] args) {
	 try {
             //https://jaxp.java.net/1.5/JAXP1.5Guide.html
             // To fix error: Caused by: 
                 //org.xml.sax.SAXParseException; systemId: file:/home/gnewton/gocode/src/github.com/gnewton/chidley/xml/MozartTrio.xml; lineNumber: 2; columnNumber: 123; External DTD: Failed to read external DTD 'partwise.dtd', because 'http' access is not allowed due to restriction set by the accessExternalDTD property.
	     System.setProperty("javax.xml.accessExternalSchema", "all");
	     System.setProperty("javax.xml.accessExternalDTD", "all");
	     System.setProperty("javax.xml.XMLConstants.ACCESS_EXTERNAL_STYLESHEET", "all");

//             System.setProperty("http.agent", "Mozilla/4.76");
               System.setProperty("http.agent", "Mozilla/5.0 (X11; Linux x86_64; rv:47.0) Gecko/20100101 Firefox/47.0");

		File file = new File("{{.SourceXMLFilename}}");
		JAXBContext jaxbContext = JAXBContext.newInstance({{.BaseXMLClassName}}.class);
 
		Unmarshaller jaxbUnmarshaller = jaxbContext.createUnmarshaller();
                {{.BaseXMLClassName}} root = ({{.BaseXMLClassName}}) jaxbUnmarshaller.unmarshal(file);

		Gson gson = new GsonBuilder().setPrettyPrinting().create();
		System.out.println(gson.toJson(root));

	  } catch (Throwable e) {
		e.printStackTrace();
                // Print classpath
                System.err.println("CLASSPATH START");
                ClassLoader cl = ClassLoader.getSystemClassLoader();
                URL[] urls = ((URLClassLoader)cl).getURLs();
                for(URL url: urls){
            	    System.err.println("\n" + url.getFile());
                }
                System.err.println("CLASSPATH END");
	  }
	}
}
`

const jaxbPackageInfoTemplage = `
@XmlSchema(
    namespace="{{.BaseNameSpace}}",
    elementFormDefault = XmlNsForm.QUALIFIED{{if .AdditionalNameSpace}},
    xmlns={
       {{range .AdditionalNameSpace}}
               @XmlNs(prefix="{{.abbr}}", namespaceURI="{{.space}}"), 
       {{end}}
   }
   {{end}}
)
package {{.PackageName}};

import javax.xml.bind.annotation.*;
`
