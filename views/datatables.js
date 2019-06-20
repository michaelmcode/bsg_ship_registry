$('#campaignTable').DataTable({
"ajax": "/api/names",

"columns": [

    {"data": "id"},
    {"data": "projectnames"},
    {"data": "universe"},
    {"data": "creationdate"},
    {"data": "image"},
]
});