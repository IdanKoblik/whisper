#!/usr/bin/env python3
import yaml
import os
import sys
import re

SWAGGER_FILE = "api/docs/swagger.yaml"
README_FILE = "README.md"

if not os.path.exists(SWAGGER_FILE):
    print(f"Swagger file not found: {SWAGGER_FILE}", file=sys.stderr)
    sys.exit(1)

if not os.path.exists(README_FILE):
    print(f"README file not found: {README_FILE}", file=sys.stderr)
    sys.exit(1)

with open(SWAGGER_FILE, "r") as f:
    swagger = yaml.safe_load(f)

paths = swagger.get("paths", {})
definitions = swagger.get("definitions", {})

def print_table(headers, rows):
    widths = [len(h) for h in headers]
    for row in rows:
        for i, cell in enumerate(row):
            widths[i] = max(widths[i], len(cell))

    header_line = "| " + " | ".join(h.ljust(widths[i]) for i, h in enumerate(headers)) + " |"
    separator = "| " + " | ".join("-" * widths[i] for i in range(len(headers))) + " |"
    table = [header_line, separator]
    for row in rows:
        table.append("| " + " | ".join(row[i].ljust(widths[i]) for i in range(len(row))) + " |")
    return "\n".join(table) + "\n"

def expand_body_param(param):
    if "$ref" in param.get("schema", {}):
        ref = param["schema"]["$ref"]
        name = ref.split("/")[-1]
        defn = definitions.get(name, {})
        props = defn.get("properties", {})
        rows = []
        for k, v in props.items():
            typ = v.get("type", "string")
            if typ == "array":
                items = v.get("items", {})
                item_type = items.get("type", "string")
                typ = f"array[{item_type}]"
            rows.append([f"`{k}`", f"`{typ}`", ""])
        return rows
    return []

api_md = "## API Reference\n\n"

for path, methods in paths.items():
    for method, info in methods.items():
        summary = info.get("summary", path)
        description = info.get("description", "")

        api_md += f"#### {summary}\n\n"
        api_md += f"{description}\n\n" if description else ""

        api_md += "```http\n"
        api_md += f"{method.upper()} {path}\n"
        api_md += "```\n\n"

        parameters = info.get("parameters", [])
        header_rows = []
        body_rows = []

        for param in parameters:
            param_in = param.get("in", "")
            if param_in == "header":
                name = f"`{param.get('name', '')}`"
                typ = f"`{param.get('type', 'string')}`"
                desc = param.get("description", "")
                if param.get("required", False):
                    desc = "**Required**. " + desc
                header_rows.append([name, typ, desc])
            elif param_in == "body":
                body_rows += expand_body_param(param)

        if header_rows:
            api_md += "\n**Header Parameters:**\n\n"
            api_md += print_table(["Parameter", "Type", "Description"], header_rows)

        if body_rows:
            api_md += "\n**Body Parameters:**\n\n"
            api_md += print_table(["Parameter", "Type", "Description"], body_rows)

        responses = info.get("responses", {})
        if responses:
            resp_rows = []
            for code, resp in responses.items():
                desc = resp.get("description", "")
                resp_rows.append([code, desc])
            api_md += "\n**Responses:**\n\n"
            api_md += print_table(["HTTP Code", "Description"], resp_rows)

        api_md += "<br>\n\n"

with open(README_FILE, "r") as f:
    readme_text = f.read()

pattern = r"(## API Reference\n)(.*?)(?=\n## |\Z)"
new_readme = re.sub(pattern, api_md.strip() + "\n", readme_text, flags=re.DOTALL)

with open(README_FILE, "w") as f:
    f.write(new_readme)

print("README.md updated successfully!")

