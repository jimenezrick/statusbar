#ifndef WM_NAME_H
#define WM_NAME_H

typedef struct xconn_t xconn_t;

xconn_t *connect_x(void);
void disconnect_x(xconn_t *xconn);
int set_screen(xconn_t *xconn);
int set_wm_name(xconn_t *xconn, const char *name);

#endif
