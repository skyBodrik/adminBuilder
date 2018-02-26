<?php

/**
  * @adminBuilder
  *     page {name: "tropic", caption: "Админка тропических Историй"},
  *     requestUrl "tropic",
  *     snippetsPath "./stories";
  */
class Story_Admin_Controller_Tropic extends Base_Controller_Admin
{

    public function indexAction()
    {
        $this->tpl = 'Tropic/index.phtml';
        $this->view->subTpl = '_main.phtml';
        $userId = $this->p('userId', $this->USER->getId());


        $game = $this->_getGame($userId);

        $this->view->storyGame = $game;

    }

    /**
     *  Получить инфу о пользователе
     *  @adminBuilder
     *      caption "Информация о пользователе",
     *      render on,
     *      in userId,
     *      out userName{snippet: "textbox", type: "string", caption: "Имя пользователя"},
     *      out userPass{snippet: "textbox", type: "string", caption: "Пароль"},
     *      linkedActions setUserName, setUserName2;
     *  @adminBuilder test;
     **/
    public function getUserInfoAction()
    {
        $this->norender();
        $userId = $this->p('userId', $this->USER->getId());
        return [
            'userName' => 'Игорь',
            'userPass' => '123',
        ];
    }

    /**
     *  Получить список пользователей
     *  @adminBuilder
     *      title "Информация о пользователе",
     *      field users{snippet: "table", type: "array", caption: "Список пользователей", fields: ["userName", "userPass"]},
     *      field userName{snippet: "textbox", type: "string", caption: "Имя пользователя"},
     *      field userPass{snippet: "textbox", type: "string", caption: "Пароль"},
     *      show users, users;
     **/
    public function listUsersAction()
    {
        return [
            users => [
                [
                    'userName' => 'Игорь',
                    'userPass' => '123',
                ],
                [
                    'userName' => 'Олег',
                    'userPass' => '204',
                ],
            ],
        ];
    }

    /**
     *  @adminBuilder render off
     *  @adminBuilder in userId, newUserName
     *  @adminBuilder out status,
     **/
    public function setUserNameAction()
    {
        $this->norender();
        $userId = $this->p('userId', $this->USER->getId());
        $questId = $this->p('newUserName', 0);
    }
}