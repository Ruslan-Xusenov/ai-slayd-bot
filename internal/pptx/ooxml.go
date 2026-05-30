package pptx

import (
	"fmt"
	"strings"
	"time"
)

// ================================================================
// [Content_Types].xml
// ================================================================

func contentTypesXML(slideCount int) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
  <Override PartName="/ppt/presProps.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presProps+xml"/>
  <Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
  <Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
  <Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
  <Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>`)
	for i := 1; i <= slideCount; i++ {
		sb.WriteString(fmt.Sprintf(`
  <Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i))
	}
	sb.WriteString("\n</Types>")
	return sb.String()
}

// ================================================================
// _rels/.rels
// ================================================================

func rootRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`
}

// ================================================================
// docProps/app.xml
// ================================================================

func appXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties">
  <Application>AI Slayd Bot</Application>
  <Slides>0</Slides>
</Properties>`
}

// ================================================================
// docProps/core.xml
// ================================================================

func coreXML(title string) string {
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
                   xmlns:dc="http://purl.org/dc/elements/1.1/"
                   xmlns:dcterms="http://purl.org/dc/terms/"
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <dc:title>%s</dc:title>
  <dc:creator>AI Slayd Bot</dc:creator>
  <dcterms:created xsi:type="dcterms:W3CDTF">%s</dcterms:created>
  <dcterms:modified xsi:type="dcterms:W3CDTF">%s</dcterms:modified>
</cp:coreProperties>`, xe(title), now, now)
}

// ================================================================
// ppt/presentation.xml
// ================================================================

func presentationXML(slideCount int) string {
	var ids strings.Builder
	for i := 1; i <= slideCount; i++ {
		ids.WriteString(fmt.Sprintf(`    <p:sldId id="%d" r:id="rId%d"/>`+"\n", 255+i, i))
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
                xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
                xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
                saveSubsetFonts="1">
  <p:sldMasterIdLst>
    <p:sldMasterId id="2147483648" r:id="rId%d"/>
  </p:sldMasterIdLst>
  <p:sldIdLst>
%s  </p:sldIdLst>
  <p:sldSz cx="%d" cy="%d" type="screen16x9"/>
  <p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`, slideCount+1, ids.String(), slideW, slideH)
}

// ================================================================
// ppt/_rels/presentation.xml.rels
// ================================================================

func presentationRelsXML(slideCount int) string {
	const (
		slideType  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
		masterType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
		propsType  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/presProps"
		themeType  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	)
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i := 1; i <= slideCount; i++ {
		sb.WriteString(fmt.Sprintf(`
  <Relationship Id="rId%d" Type="%s" Target="slides/slide%d.xml"/>`, i, slideType, i))
	}
	masterID := slideCount + 1
	sb.WriteString(fmt.Sprintf(`
  <Relationship Id="rId%d" Type="%s" Target="slideMasters/slideMaster1.xml"/>
  <Relationship Id="rId%d" Type="%s" Target="presProps.xml"/>
  <Relationship Id="rId%d" Type="%s" Target="theme/theme1.xml"/>
</Relationships>`, masterID, masterType, masterID+1, propsType, masterID+2, themeType))
	return sb.String()
}

// ================================================================
// ppt/presProps.xml
// ================================================================

func presPropsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentationPr xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
                  xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
                  xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"/>`
}

// ================================================================
// ppt/theme/theme1.xml  (minimal, generic)
// ================================================================

func theme1XML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="AIBot">
  <a:themeElements>
    <a:clrScheme name="AIBot">
      <a:dk1><a:sysClr lastClr="000000" val="windowText"/></a:dk1>
      <a:lt1><a:sysClr lastClr="FFFFFF" val="window"/></a:lt1>
      <a:dk2><a:srgbClr val="44546A"/></a:dk2>
      <a:lt2><a:srgbClr val="E7E6E6"/></a:lt2>
      <a:accent1><a:srgbClr val="4472C4"/></a:accent1>
      <a:accent2><a:srgbClr val="ED7D31"/></a:accent2>
      <a:accent3><a:srgbClr val="A9D18E"/></a:accent3>
      <a:accent4><a:srgbClr val="FFC000"/></a:accent4>
      <a:accent5><a:srgbClr val="5A96C7"/></a:accent5>
      <a:accent6><a:srgbClr val="70AD47"/></a:accent6>
      <a:hlink><a:srgbClr val="0563C1"/></a:hlink>
      <a:folHlink><a:srgbClr val="954F72"/></a:folHlink>
    </a:clrScheme>
    <a:fontScheme name="AIBot">
      <a:majorFont>
        <a:latin typeface="Calibri Light"/>
        <a:ea typeface=""/>
        <a:cs typeface=""/>
      </a:majorFont>
      <a:minorFont>
        <a:latin typeface="Calibri"/>
        <a:ea typeface=""/>
        <a:cs typeface=""/>
      </a:minorFont>
    </a:fontScheme>
    <a:fmtScheme name="AIBot">
      <a:fillStyleLst>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
      </a:fillStyleLst>
      <a:lnStyleLst>
        <a:ln w="6350"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
        <a:ln w="12700"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
        <a:ln w="19050"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill></a:ln>
      </a:lnStyleLst>
      <a:effectStyleLst>
        <a:effectStyle><a:effectLst/></a:effectStyle>
        <a:effectStyle><a:effectLst/></a:effectStyle>
        <a:effectStyle><a:effectLst/></a:effectStyle>
      </a:effectStyleLst>
      <a:bgFillStyleLst>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
        <a:solidFill><a:schemeClr val="phClr"/></a:solidFill>
      </a:bgFillStyleLst>
    </a:fmtScheme>
  </a:themeElements>
</a:theme>`
}

// ================================================================
// ppt/slideMasters/slideMaster1.xml
// ================================================================

func slideMasterXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
             xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
             xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:bg>
      <p:bgRef idx="1001"><a:schemeClr val="bg1"/></p:bgRef>
    </p:bg>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="0" cy="0"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="0" cy="0"/>
        </a:xfrm>
      </p:grpSpPr>
    </p:spTree>
  </p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2"
            accent1="accent1" accent2="accent2" accent3="accent3"
            accent4="accent4" accent5="accent5" accent6="accent6"
            hlink="hlink" folHlink="folHlink"/>
  <p:sldLayoutIdLst>
    <p:sldLayoutId id="2147483648" r:id="rId1"/>
  </p:sldLayoutIdLst>
  <p:txStyles>
    <p:titleStyle>
      <a:lvl1pPr algn="l">
        <a:defRPr sz="3600" b="1">
          <a:solidFill><a:schemeClr val="tx1"/></a:solidFill>
          <a:latin typeface="+mj-lt"/>
        </a:defRPr>
      </a:lvl1pPr>
    </p:titleStyle>
    <p:bodyStyle>
      <a:lvl1pPr>
        <a:defRPr sz="2000">
          <a:solidFill><a:schemeClr val="tx1"/></a:solidFill>
          <a:latin typeface="+mn-lt"/>
        </a:defRPr>
      </a:lvl1pPr>
    </p:bodyStyle>
    <p:otherStyle>
      <a:defPPr><a:defRPr lang="uz-UZ"/></a:defPPr>
    </p:otherStyle>
  </p:txStyles>
</p:sldMaster>`
}

func slideMasterRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme1.xml"/>
</Relationships>`
}

// ================================================================
// ppt/slideLayouts/slideLayout1.xml  (blank layout)
// ================================================================

func slideLayoutXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
             xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
             xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"
             type="blank" preserve="1">
  <p:cSld name="Blank">
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="0" cy="0"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="0" cy="0"/>
        </a:xfrm>
      </p:grpSpPr>
    </p:spTree>
  </p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sldLayout>`
}

func slideLayoutRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="../slideMasters/slideMaster1.xml"/>
</Relationships>`
}

// ================================================================
// ppt/slides/_rels/slideN.xml.rels
// ================================================================

func slideRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout1.xml"/>
</Relationships>`
}

// ================================================================
// Slayd XML generatorlar
// ================================================================

// solidFill — inline rang to'ldirish
func solidFill(hex string) string {
	return fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, hex)
}

// sp — to'rtburchak matn qutisi
func sp(id int, name, x, y, cx, cy string, body string) string {
	return fmt.Sprintf(`<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="%d" name="%s"/>
    <p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm><a:off x="%s" y="%s"/><a:ext cx="%s" cy="%s"/></a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
    <a:noFill/>
  </p:spPr>
  %s
</p:sp>`, id, xe(name), x, y, cx, cy, body)
}

// txBody — matn qutisi ichidagi txBody
func txBody(anchor, paragraphs string) string {
	return fmt.Sprintf(`<p:txBody>
    <a:bodyPr wrap="square" rtlCol="0" anchor="%s"/>
    <a:lstStyle/>
    %s
  </p:txBody>`, anchor, paragraphs)
}

// para — bitta paragraf
func para(align, color string, sz int, bold bool, text string) string {
	b := "0"
	if bold {
		b = "1"
	}
	return fmt.Sprintf(`<a:p>
      <a:pPr algn="%s"/>
      <a:r>
        <a:rPr lang="uz-UZ" sz="%d" b="%s" dirty="0">
          %s
          <a:latin typeface="Calibri"/>
        </a:rPr>
        <a:t>%s</a:t>
      </a:r>
    </a:p>`, align, sz, b, solidFill(color), xe(text))
}

// rect — to'ldirilgan to'rtburchak (dekoratsiya uchun)
func rect(x, y, cx, cy int, color string) string {
	return fmt.Sprintf(`<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="99" name="rect"/>
    <p:cNvSpPr/>
    <p:nvPr/>
  </p:nvSpPr>
  <p:spPr>
    <a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>
    <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
    %s
    <a:ln><a:noFill/></a:ln>
  </p:spPr>
  <p:txBody><a:bodyPr/><a:lstStyle/><a:p/></p:txBody>
</p:sp>`, x, y, cx, cy, solidFill(color))
}

// bgShape — fon rangi (to'liq slayd)
func bgXML(t Theme) string {
	if t.BgGradient != "" {
		return fmt.Sprintf(`<p:bg>
  <p:bgPr>
    <a:gradFill>
      <a:gsLst>
        <a:gs pos="0"><a:srgbClr val="%s"/></a:gs>
        <a:gs pos="100000"><a:srgbClr val="%s"/></a:gs>
      </a:gsLst>
      <a:lin ang="5400000" scaled="0"/>
    </a:gradFill>
    <a:effectLst/>
  </p:bgPr>
</p:bg>`, t.Background, t.BgGradient)
	}
	return fmt.Sprintf(`<p:bg>
  <p:bgPr>
    %s
    <a:effectLst/>
  </p:bgPr>
</p:bg>`, solidFill(t.Background))
}

// sldHeader — har bir slayd uchun umumiy wrapper
func sldWrap(bg, shapes string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
       xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
       xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    %s
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="%d" cy="%d"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="%d" cy="%d"/>
        </a:xfrm>
      </p:grpSpPr>
      %s
    </p:spTree>
  </p:cSld>
  <p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr>
</p:sld>`, bg, slideW, slideH, slideW, slideH, shapes)
}

// ================================================================
// Sarlavha slayd
// ================================================================

func titleSlideXML(t Theme, title, subtitle string) string {
	// Tepa chiziq
	topBar := rect(0, 0, slideW, 18000, t.HeaderBar)

	// Asosiy sarlavha
	titleBox := sp(2, "Title",
		"914400", "2286000", "10363200", "1600000",
		txBody("ctr", para("ctr", t.Title, 5400, true, title)),
	)

	// Kichik sarlavha
	subBox := sp(3, "Subtitle",
		"914400", "3886000", "10363200", "914400",
		txBody("t", para("ctr", t.Subtitle, 2200, false, subtitle)),
	)

	// Accent chiziq
	accentLine := rect(4572000, 4914000, 3048000, 27000, t.Accent)

	// Pastki chiziq
	bottomBar := rect(0, slideH-18000, slideW, 18000, t.HeaderBar)

	shapes := topBar + "\n" + titleBox + "\n" + subBox + "\n" + accentLine + "\n" + bottomBar
	return sldWrap(bgXML(t), shapes)
}

// ================================================================
// Kontent slayd
// ================================================================

func contentSlideXML(t Theme, title string, bullets []string, presTitle string, idx, total int) string {
	topBar := rect(0, 0, slideW, 18000, t.HeaderBar)
	sideBar := rect(457200, 457200, 72000, 685800, t.Accent)

	titleBox := sp(2, "Title",
		"600000", "380000", "11135000", "800000",
		txBody("b", para("l", t.Title, 3200, true, title)),
	)

	// Accent chiziq (sarlavha ostida)
	accentLine := rect(457200, 1260000, 2743200, 18000, t.Accent)

	// Bullet paragraflar
	var bulletParas strings.Builder
	for _, b := range bullets {
		if strings.TrimSpace(b) == "" {
			continue
		}
		bulletParas.WriteString(fmt.Sprintf(`<a:p>
      <a:pPr marL="342900" indent="-342900">
        <a:buChar char="▪"/>
      </a:pPr>
      <a:r>
        <a:rPr lang="uz-UZ" sz="1900" b="0" dirty="0">
          %s
          <a:latin typeface="Calibri"/>
        </a:rPr>
        <a:t>%s</a:t>
      </a:r>
    </a:p>`, solidFill(t.Body), xe(strings.TrimSpace(b))))
	}

	contentBox := sp(3, "Content",
		"457200", "1371600", "11277600", "4800000",
		fmt.Sprintf(`<p:txBody>
    <a:bodyPr wrap="square" rtlCol="0" anchor="t"/>
    <a:lstStyle/>
    %s
  </p:txBody>`, bulletParas.String()),
	)

	// Footer
	footerText := fmt.Sprintf("%s  |  %d / %d", trunc(presTitle, 50), idx, total)
	footerBox := sp(4, "Footer",
		"457200", "6400000", "11277600", "300000",
		txBody("b", para("l", t.Footer, 1200, false, footerText)),
	)

	bottomBar := rect(0, slideH-18000, slideW, 18000, t.HeaderBar)

	shapes := topBar + "\n" + sideBar + "\n" + titleBox + "\n" + accentLine + "\n" + contentBox + "\n" + footerBox + "\n" + bottomBar
	return sldWrap(bgXML(t), shapes)
}

// ================================================================
// Yakuniy slayd
// ================================================================

func endSlideXML(t Theme, presTitle string) string {
	topBar := rect(0, 0, slideW, 18000, t.HeaderBar)

	thankBox := sp(2, "ThankYou",
		"914400", "2057400", "10363200", "1371600",
		txBody("ctr", para("ctr", t.Title, 6000, true, "Thank you!")),
	)

	subBox := sp(3, "PresTitle",
		"914400", "3657600", "10363200", "685800",
		txBody("t", para("ctr", t.Subtitle, 2000, false, presTitle)),
	)

	accentLine := rect(4572000, 4457200, 3048000, 27000, t.Accent)
	bottomBar := rect(0, slideH-18000, slideW, 18000, t.HeaderBar)

	shapes := topBar + "\n" + thankBox + "\n" + subBox + "\n" + accentLine + "\n" + bottomBar
	return sldWrap(bgXML(t), shapes)
}
